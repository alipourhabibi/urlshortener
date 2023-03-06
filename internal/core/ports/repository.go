package ports

import (
	aggregate "github.com/alipourhabibi/urlshortener/internal/core/aggergate"
)

type URLRepository interface {
	Add(aggregate.URL) (string, error)
	Exists(string) (bool, error)
	Get(string) (aggregate.URL, error)
}
