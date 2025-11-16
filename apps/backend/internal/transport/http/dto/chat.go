package dto

import "time"

type HealthResponse struct {
	Status string `json:"status"`
}

type ChatResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ChatsListResponse struct {
	Chats []ChatResponse `json:"chats"`
}

type CreateChatRequest struct {
	Title        string  `json:"title"`
	ScenarioCode *string `json:"scenario_code,omitempty"`
}

type CreateChatResponse struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type MessageResponse struct {
	ID        string    `json:"id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type MessagesListResponse struct {
	Messages []MessageResponse `json:"messages"`
}

type SendMessageRequest struct {
	Content      string   `json:"content"`
	DocumentIDs  []string `json:"document_ids,omitempty"`
	ScenarioCode *string  `json:"scenario_code,omitempty"`
}

type SendMessageResponse struct {
	Message MessageResponse `json:"message"`
}
