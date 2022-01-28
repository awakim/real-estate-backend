package cache

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type Cache interface {
	SetRefreshToken(ctx context.Context, username string, tokenID string, expiration time.Duration) error
	DeleteRefreshToken(ctx context.Context, username string, tokenID string) error
}

type RedisStore struct {
	client *redis.Client
}

func NewCache(client *redis.Client) Cache {
	return &RedisStore{client: client}
}
