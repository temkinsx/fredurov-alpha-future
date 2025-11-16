package llm

import (
	"alpha_future_fredurov/apps/backend/internal/domain"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// DocumentTextGetter - интерфейс для получения текста документа
// Это может быть реализовано через отдельный сервис или метод в DocumentRepo
type DocumentTextGetter interface {
	GetDocumentText(ctx context.Context, docID uuid.UUID) (string, error)
}

// Reply обрабатывает сообщение пользователя и возвращает ответ от LLM
func (s *Service) Reply(
	ctx context.Context,
	chatID uuid.UUID,
	userID uuid.UUID,
	userText string,
	documentIDs []uuid.UUID,
	scenarioCode *string,
	docTextGetter DocumentTextGetter,
) (*domain.Message, error) {
	if userText == "" {
		return nil, errors.New("user text cannot be empty")
	}

	// 1. Проверяем чат и права доступа
	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat: %w", err)
	}

	if chat.UserID != userID {
		return nil, errors.New("access denied: chat belongs to different user")
	}

	// 2. Создаём сообщение пользователя
	userMsg := &domain.Message{
		ID:      uuid.New(),
		ChatID:  chatID,
		Role:    string(domain.RoleUser),
		Content: userText,
	}

	if err := s.msgRepo.Append(ctx, userMsg); err != nil {
		return nil, fmt.Errorf("failed to save user message: %w", err)
	}

	// 3. Получаем историю сообщений
	rawMsgHistory, err := s.msgRepo.GetLastN(ctx, chatID, 50) // берем больше, потом обрежем по лимитам
	if err != nil {
		return nil, fmt.Errorf("failed to get message history: %w", err)
	}

	// Реверсируем порядок (GetLastN возвращает от новых к старым, нужны от старых к новым)
	var msgHistory strings.Builder
	for i := len(rawMsgHistory) - 1; i >= 0; i-- {
		msg := rawMsgHistory[i]
		msgHistory.WriteString(fmt.Sprintf("%s: %s\n", msg.Role, msg.Content))
	}

	// 4. Получаем текст документов
	var documentsText strings.Builder
	if docTextGetter != nil && len(documentIDs) > 0 {
		for _, docID := range documentIDs {
			text, err := docTextGetter.GetDocumentText(ctx, docID)
			if err != nil {
				// Логируем ошибку, но продолжаем
				continue
			}
			if text != "" {
				documentsText.WriteString(text)
				documentsText.WriteString("\n\n")
			}
		}
	}

	// 5. Выбираем системный промпт (по сценарию или дефолтный)
	sysPrompt := s.getSystemPrompt(scenarioCode)

	// 6. Собираем промпт
	prompt := s.buildPrompt(
		sysPrompt,
		msgHistory.String(),
		documentsText.String(),
		userText,
	)

	// 7. Вызываем LLM
	startTime := time.Now()
	llmResponse, err := s.llm.Generate(ctx, []byte(prompt))
	if err != nil {
		return nil, fmt.Errorf("failed to generate LLM response: %w", err)
	}
	latencyMs := time.Since(startTime).Milliseconds()

	// 8. Создаём сообщение ассистента
	assistantMsg := &domain.Message{
		ID:        uuid.New(),
		ChatID:    chatID,
		Role:      string(domain.RoleAssistant),
		Content:   llmResponse,
		LatencyMs: &latencyMs,
	}

	if err := s.msgRepo.Append(ctx, assistantMsg); err != nil {
		return nil, fmt.Errorf("failed to save assistant message: %w", err)
	}

	// 9. Обновляем время последнего сообщения в чате
	now := time.Now()
	if err := s.chatRepo.Touch(ctx, chatID, now); err != nil {
		// Логируем, но не возвращаем ошибку
	}

	return assistantMsg, nil
}

// getSystemPrompt возвращает системный промпт по коду сценария
func (s *Service) getSystemPrompt(scenarioCode *string) string {
	if scenarioCode == nil {
		return "" // будет использован дефолтный из buildPrompt
	}

	// TODO: реализовать маппинг сценариев на промпты
	// Пока возвращаем пустую строку для использования дефолтного
	switch *scenarioCode {
	case "contract_helper":
		return "Ты — помощник по анализу договоров. Твоя задача: объяснить условия договора простым языком, выделить риски и важные моменты, помочь подготовить формулировки для переговоров."
	case "marketing":
		return "Ты — маркетинговый ассистент для микробизнеса. Помогай создавать посты, придумывать акции, писать тексты для соцсетей и рекламы."
	default:
		return "" // дефолтный промпт
	}
}
