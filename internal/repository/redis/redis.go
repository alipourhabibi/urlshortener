package redis

import (
	"time"

	"github.com/alipourhabibi/urlshortener/config"
	"github.com/go-redis/redis"
)

type RedisRepository struct {
	db *redis.Client
}

func New(confs config.Redis) (*RedisRepository, error) {
	db := redis.NewClient(&redis.Options{
		Addr:     confs.Host + ":" + confs.Port,
		Password: confs.Password,
		DB:       confs.DBName,
	})
	err := db.Ping().Err()
	if err != nil {
		return nil, err
	}
	return &RedisRepository{
		db: db,
	}, nil
}

func (r *RedisRepository) Set(key string, value interface{}, t time.Duration) (interface{}, error) {
	v := r.db.Set(key, value, t)
	return v.Val(), v.Err()
}

func (r *RedisRepository) Get(key string) (interface{}, error) {
	v := r.db.Get(key)
	return v.Val(), v.Err()
}
