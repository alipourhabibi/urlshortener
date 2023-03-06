package entity

import "github.com/google/uuid"

type Shortened struct {
	ID           uuid.UUID
	OriginalURL  string
	ShortenedURL string
}
