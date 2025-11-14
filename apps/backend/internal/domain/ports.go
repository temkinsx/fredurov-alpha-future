package domain

import (
	"context"

	"github.com/google/uuid"
)

type ChatRepo interface {
	// Create - сохранить новый чат в БД
	Create(ctx context.Context, chat *Chat) error
	// Get - получить чат по ID
	Get(ctx context.Context, chatID uuid.UUID) (*Chat, error)
	// ListByUser - получить все чаты пользователя по его ID
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*Chat, error)
	// Update - Обновить существующий чат
	Update(ctx context.Context, chat *Chat) error
	// Delete - удалить чат из БД по ID
	Delete(ctx context.Context, chatID uuid.UUID) error
}

type MessageRepo interface {
	// Append — сохранить новое сообщение в конец чата.
	// Предполагаем, что msg.ChatID уже заполнен.
	Append(ctx context.Context, msg *Message) error
	// GetLastN — взять последние n сообщений чата.
	// Порядок - от старых к новым
	GetLastN(ctx context.Context, chatID uuid.UUID, n int) ([]*Message, error)
	// ListByChat — история чата для UI с пагинацией.
	// Порядок - от старых к новым
	ListByChat(ctx context.Context, chatID uuid.UUID, limit, offset int) ([]*Message, error)
}

type LLM interface {
	// Generate - отправить запрос к LLM. Возвращает ответ в виде string
	Generate(ctx context.Context, prompt []byte) (string, error)
}
