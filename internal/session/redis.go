package session

import (
	"context"

	"github.com/go-redis/redis"
)

// Config - db
type Config struct {
	Url string
}

// RedisClient
type RedisClient struct {
	ctx    context.Context
	client *redis.Client
}

// NewRedisClient - constructor
func NewRedisClient(cfg Config) (*RedisClient, error) {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Url,
	})

	return &RedisClient{
		ctx:    ctx,
		client: client,
	}, nil
}
