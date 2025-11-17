package domain

import "time"

type Document struct {
	ID         string
	UserID     string
	ChatID     *string
	Name       string
	MimeType   string
	StorageURL string
	SizeBytes  int64

	CreatedAt time.Time
}

