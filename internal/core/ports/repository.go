package ports

import (
	"time"

	aggregate "github.com/alipourhabibi/urlshortener/internal/core/aggergate"
)

type URLRepository interface {
	Add(aggregate.URL) (string, error)
	Exists(string) (bool, error)
	Get(string) (aggregate.URL, error)
}

type CacheRepository interface {
	Set(string, interface{}, time.Duration) (interface{}, error)
	Get(string) (interface{}, error)
}
