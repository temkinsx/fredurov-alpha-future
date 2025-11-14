package domain

import (
	"time"

	"github.com/google/uuid"
)

type ChatStatus string

const (
	ChatActive   ChatStatus = "active"
	ChatArchived ChatStatus = "archived"
)

type Chat struct {
	ID       uuid.UUID         `json:"id"`
	Title    *string           `json:"title,omitempty"`
	UserID   uuid.UUID         `json:"user_id,omitempty"`
	Model    *string           `json:"model,omitempty"` // log
	Summary  *string           `json:"summary,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Status   ChatStatus        `json:"status"`

	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	LastMessageAt time.Time `json:"last_message_at"`
}
