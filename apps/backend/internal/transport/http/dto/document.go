package dto

import "time"

type DocumentResponse struct {
	ID        string    `json:"id"`
	FileName  string    `json:"file_name"`
	ChatID    *string   `json:"chat_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type DocumentsListResponse struct {
	Documents []DocumentResponse `json:"documents"`
}

type DocumentDetailResponse struct {
	ID          string    `json:"id"`
	FileName    string    `json:"file_name"`
	MimeType    string    `json:"mime_type"`
	ChatID      *string   `json:"chat_id,omitempty"`
	SizeBytes   int64     `json:"size_bytes"`
	CreatedAt   time.Time `json:"created_at"`
	TextExcerpt string    `json:"text_excerpt"`
}

type CreateDocumentResponse struct {
	Document DocumentResponse `json:"document"`
}
