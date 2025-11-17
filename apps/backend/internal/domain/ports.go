package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ChatRepo interface {
	// Create - сохранить новый чат в БД
	Create(ctx context.Context, chat *Chat) error
	// Get - получить чат по ID
	GetByID(ctx context.Context, chatID uuid.UUID) (*Chat, error)
	// ListByUser - получить все чаты пользователя по его ID
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*Chat, error)
	// Update - Обновить существующий чат
	UpdateTitle(ctx context.Context, chat *Chat) error
	//Touch - обновление времени последнего сообщения в чате (last_message_at)
	Touch(ctx context.Context, chatID uuid.UUID, t time.Time) error
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
	// Для хендлеров стоит в main.go создать новый сервис ser
	Generate(ctx context.Context, prompt []byte) (string, error)
}

type UserRepo interface {
	// GetByEmail - получить пользователя по email
	GetByEmail(ctx context.Context, email string) (*User, error)
	// UpdateLastLogin - обновить время последнего входа
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
}
