package cache

import (
	"context"
	"time"

	"github.com/awakim/immoblock-backend/token"
	"github.com/go-redis/redis/v8"
)

type Cache interface {
	// DeleteRefreshtoken deletes old refresh tokens
	DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error
	// SetTokenData sets the access and refresh tokens related data into Redis cache.
	SetTokenData(ctx context.Context, accessToken token.Payload, atd time.Duration, refreshToken token.Payload, rtd time.Duration) error
	// LogoutUser deletes access and refresh tokens from Cache and revoke access and refresh tokens by setting them in Cache.
	LogoutUser(ctx context.Context, accessToken token.Payload, refreshToken token.Payload) error
	// IsRevoked checks whether a token is revoked by checking the according key `rev:{{userID}}:{{tokenID}}`.
	// If the key present, the token is revoked, else perhaps a server error and finally if none of the
	// previous then token is not revoked.
	IsRevoked(ctx context.Context, token token.Payload) (bool, error)
	// IsRateLimited checks whether a user has surpassed the limit of login or refresh routes.
	// The rate limit is imposed as 3 requests per IP per Identifier (Email for login or UserID for refresh)
	// per quarter hour.
	IsRateLimited(ctx context.Context, identifier string) (bool, error)
}

type RedisStore struct {
	Client *redis.Client
}

func NewCache(client *redis.Client) Cache {
	return &RedisStore{
		Client: client,
	}
}
