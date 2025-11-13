package domain

import (
	"context"

	"github.com/google/uuid"
)

type ChatRepo interface {
	Create(ctx context.Context, chat *Chat) error
	Get(ctx context.Context, chatID uuid.UUID) (*Chat, error)
	ListByUser(ctx context.Context, userID uuid.UUID) ([]*Chat, error)
	Update(ctx context.Context, chat *Chat) error
	Delete(ctx context.Context, chatID uuid.UUID) error
}

type MessageRepo interface {
	Append(ctx context.Context, msg *Message) error
	GetLastN(ctx context.Context, chatID uuid.UUID, n int) ([]*Message, error)
}

type LLM interface {
	Generate(ctx context.Context, chunks []byte) (string, error)
}
