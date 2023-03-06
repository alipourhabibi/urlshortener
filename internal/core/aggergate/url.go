package aggergate

import (
	"fmt"

	"github.com/alipourhabibi/urlshortener/internal/core/entity"
	"github.com/google/uuid"
)

var (
	ErrInvalidURL    = fmt.Errorf("Invalid URL")
	ErrShorteningURL = fmt.Errorf("Error while shortening URL")
)

type URL struct {
	shortened *entity.Shortened
	user      *entity.User
}

func NewURL(original string, user *entity.User) (URL, error) {
	if original == "" {
		return URL{}, ErrInvalidURL
	}
	shortened := &entity.Shortened{
		ID:           uuid.New(),
		OriginalURL:  original,
		ShortenedURL: "",
	}
	return URL{
		shortened: shortened,
		user:      user,
	}, nil
}

// get id of url
func (u *URL) GetID() uuid.UUID {
	return u.shortened.ID
}

// gets the original url
func (u *URL) GetOriginal() string {
	return u.shortened.OriginalURL
}

// gets the original url
func (u *URL) GetShortened() string {
	return u.shortened.ShortenedURL
}

// set id of url
func (u *URL) SetID(id uuid.UUID) {
	u.shortened.ID = id
}

// set original url
func (u *URL) SetOriginal(original string) {
	u.shortened.OriginalURL = original
}

// set shortened url
func (u *URL) SetShortened(shortened string) {
	u.shortened.ShortenedURL = shortened
}
