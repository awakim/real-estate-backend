package cache

import (
	"context"
	"time"

	"github.com/awakim/immoblock-backend/token"
	"github.com/go-redis/redis/v8"
)

type Cache interface {
	// SetRefreshToken(ctx context.Context, userID string, tokenID string, expiration time.Duration) error
	DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error
	// LPushAccessToken(ctx context.Context, userID string, tokenID string) error
	// SetTokensInfos(ctx context.Context, userID string, atID string, rtID string) error
	SetTokenData(ctx context.Context, accessToken token.Payload, atd time.Duration, refreshToken token.Payload, rtd time.Duration) error
	GetTokenData(ctx context.Context, refreshToken token.Payload) error
}

type RedisStore struct {
	Client *redis.Client
}

func NewCache(client *redis.Client) Cache {
	return &RedisStore{
		Client: client,
	}
}
