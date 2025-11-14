package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleSystem    Role = "system" // используем для системных сообщений (например установить поведение)
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Message struct {
	ID      uuid.UUID `json:"id"`
	ChatID  uuid.UUID `json:"conversation_id"`
	Role    string    `json:"role"`
	Content string    `json:"content"`

	CreatedAt time.Time `json:"created_at"`
	LatencyMs *int64
	Truncated bool // флаг, что часть запроса/документа была обрезана по лимиту
}

func (m *Message) String() string {
	return fmt.Sprintf("%s: %s", m.Role, m.Content)
}
