package cache

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// SetRefreshToken stores a refresh Token with an expiry time
func (cache *RedisStore) SetRefreshToken(ctx context.Context, username string, tokenID string, expiration time.Duration) error {
	if tokenID == "" {
		return errors.New("tokenID cannot be equal to empty string")
	}
	if username == "" {
		return errors.New("username cannot be equal to empty string")
	}
	key := fmt.Sprintf("%s:%s", username, tokenID)
	if err := cache.client.Set(ctx, key, 0, expiration).Err(); err != nil {
		return err
	}
	return nil
}

// DeleteRefreshtoken deletes old refresh tokens
func (cache *RedisStore) DeleteRefreshToken(ctx context.Context, username string, tokenID string) error {
	key := fmt.Sprintf("%s:%s", username, tokenID)

	response := cache.client.Del(ctx, key)
	if err := response.Err(); err != nil {
		return err
	}

	if response.Val() < 1 {
		return fmt.Errorf("invalid Refresh token: %s:%s does not exist", username, tokenID)
	}

	return nil
}
