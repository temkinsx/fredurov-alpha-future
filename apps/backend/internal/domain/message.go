package domain

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID      uuid.UUID `json:"id"`
	ChatID  uuid.UUID `json:"conversation_id"`
	Role    string
	Content string `json:"content"`

	CreatedAt time.Time `json:"created_at"`

	//	ToolCalls  []ToolCall  `json:"tool_calls,omitempty"`
	//	ToolResult *ToolResult `json:"tool_result,omitempty"`
	//	Citations []Citation `json:"citations,omitempty"`
	//Usage *Usage `json:"usage,omitempty"`
}
