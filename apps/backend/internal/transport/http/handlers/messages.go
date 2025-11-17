package handlers

import (
	"backend/internal/domain"
	"backend/internal/transport/http/dto"
	"backend/internal/usecase/llm"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type MessagesHandler struct {
	msgRepo       domain.MessageRepo
	chatRepo      domain.ChatRepo
	llmService    *llm.Service
	docTextGetter llm.DocumentTextGetter
}

func NewMessagesHandler(
	msgRepo domain.MessageRepo,
	chatRepo domain.ChatRepo,
	llmService *llm.Service,
	docTextGetter llm.DocumentTextGetter,
) *MessagesHandler {
	return &MessagesHandler{
		msgRepo:       msgRepo,
		chatRepo:      chatRepo,
		llmService:    llmService,
		docTextGetter: docTextGetter,
	}
}

// GetMessages возвращает историю сообщений чата
func (h *MessagesHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	chatIDStr := chi.URLParam(r, "chat_id")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		http.Error(w, "invalid chat_id", http.StatusBadRequest)
		return
	}

	// Проверяем права доступа
	chat, err := h.chatRepo.GetByID(r.Context(), chatID)
	if err != nil {
		http.Error(w, "chat not found", http.StatusNotFound)
		return
	}

	if chat.UserID != userID {
		http.Error(w, "access denied", http.StatusForbidden)
		return
	}

	messages, err := h.msgRepo.ListByChat(r.Context(), chatID, 100, 0)
	if err != nil {
		http.Error(w, "failed to get messages", http.StatusInternalServerError)
		return
	}

	response := dto.MessagesListResponse{
		Messages: make([]dto.MessageResponse, len(messages)),
	}

	for i, msg := range messages {
		response.Messages[i] = dto.MessageResponse{
			ID:        msg.ID.String(),
			Role:      msg.Role,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SendMessage обрабатывает сообщение пользователя и возвращает ответ от LLM
func (h *MessagesHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	chatIDStr := chi.URLParam(r, "chat_id")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		http.Error(w, "invalid chat_id", http.StatusBadRequest)
		return
	}

	var req dto.SendMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Content) == "" {
		http.Error(w, "content is required", http.StatusBadRequest)
		return
	}

	// Парсим document IDs
	var documentIDs []uuid.UUID
	for _, docIDStr := range req.DocumentIDs {
		docID, err := uuid.Parse(docIDStr)
		if err != nil {
			continue // пропускаем невалидные ID
		}
		documentIDs = append(documentIDs, docID)
	}

	// Вызываем LLM сервис
	assistantMsg, err := h.llmService.Reply(
		r.Context(),
		chatID,
		userID,
		req.Content,
		documentIDs,
		req.ScenarioCode,
		h.docTextGetter,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := dto.SendMessageResponse{
		Message: dto.MessageResponse{
			ID:        assistantMsg.ID.String(),
			Role:      assistantMsg.Role,
			Content:   assistantMsg.Content,
			CreatedAt: assistantMsg.CreatedAt,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
