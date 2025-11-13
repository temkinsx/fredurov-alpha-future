package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

type ChatStatus string

const (
	ChatActive   ChatStatus = "active"
	ChatArchived ChatStatus = "archived"
)

type Chat struct {
	ID       uuid.UUID         `json:"id"`
	Tenant   string            `json:"tenant"`
	Title    *string           `json:"title,omitempty"`
	UserID   *string           `json:"user_id,omitempty"`
	Model    *string           `json:"model,omitempty"` // модель по умолчанию
	Summary  *string           `json:"summary,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Status   ChatStatus        `json:"status"`

	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	LastMessageAt time.Time `json:"last_message_at"`
}
