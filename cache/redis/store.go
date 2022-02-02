package cache

import (
	"context"
	"time"

	"github.com/awakim/immoblock-backend/token"
	"github.com/go-redis/redis/v8"
)

type Cache interface {
	DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error
	SetTokenData(ctx context.Context, accessToken token.Payload, atd time.Duration, refreshToken token.Payload, rtd time.Duration) error
	LogoutUser(ctx context.Context, accessToken token.Payload, refreshToken token.Payload) error
	IsRevoked(ctx context.Context, token token.Payload) (bool, error)
}

type RedisStore struct {
	Client *redis.Client
}

func NewCache(client *redis.Client) Cache {
	return &RedisStore{
		Client: client,
	}
}
