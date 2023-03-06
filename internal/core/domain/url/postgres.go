package url

import (
	aggregate "github.com/alipourhabibi/urlshortener/internal/core/aggergate"
	"github.com/google/uuid"
)

// postgresURL is an internal type that is used to store a URLAggregate
type PostgresURL struct {
	ID        uuid.UUID
	Original  string
	Shortened string
}

// NewFromURL - an aggregate and converts it to the postgresURL
func NewFromURL(u aggregate.URL) PostgresURL {
	return PostgresURL{
		ID:        u.GetID(),
		Original:  u.GetOriginal(),
		Shortened: u.GetShortened(),
	}
}

// ToAggregate - converts to aggregate.URL
func (p PostgresURL) ToAggregate() aggregate.URL {
	u := aggregate.URL{}

	u.SetID(p.ID)
	u.SetOriginal(p.Original)
	u.SetShortened(p.Shortened)

	return u
}
