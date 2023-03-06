package ports

import (
	aggregate "github.com/alipourhabibi/urlshortener/internal/core/aggergate"
	"github.com/google/uuid"
)

type URLRepository interface {
	Add(aggregate.URL) error
	Get(uuid.UUID) (aggregate.URL, error)
}
