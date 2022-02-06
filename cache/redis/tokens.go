package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/awakim/immoblock-backend/token"
	"github.com/go-redis/redis/v8"
)

// DeleteRefreshtoken deletes old refresh tokens
func (cache *RedisStore) DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error {
	key := fmt.Sprintf("rt:%s:%s", userID, tokenID)
	response := cache.Client.Del(ctx, key)
	if err := response.Err(); err != nil {
		return err
	}
	if response.Val() < 1 {
		return redis.Nil
	}
	return nil
}

// SetTokenData sets the access and refresh tokens related data into Redis cache.
func (cache *RedisStore) SetTokenData(ctx context.Context, accessToken token.Payload, atd time.Duration, refreshToken token.Payload, rtd time.Duration) error {
	pipe := cache.Client.TxPipeline()

	if accessToken.UserID.String() == "" || accessToken.ID.String() == "" {
		return errors.New("invalid access token data to set in cache")
	}
	at := fmt.Sprintf("at:%s:%s", accessToken.UserID.String(), accessToken.ID.String())
	pipe.SetEX(ctx, at, 1, atd)

	if refreshToken.UserID.String() == "" || refreshToken.ID.String() == "" {
		return errors.New("invalid refresh token data to set in cache")
	}
	rt := fmt.Sprintf("rt:%s:%s", refreshToken.UserID.String(), refreshToken.ID.String())
	pipe.SetEX(ctx, rt, 1, rtd)
	atlKey := fmt.Sprintf("atl:%s", accessToken.UserID.String())
	atlValue := []interface{}{accessToken.ID.String()}
	pipe.LPush(ctx, atlKey, atlValue...)
	pipe.Expire(ctx, atlKey, atd)

	rtlKey := fmt.Sprintf("rtl:%s", refreshToken.UserID.String())
	rtlValue := []interface{}{refreshToken.ID.String()}
	pipe.LPush(ctx, rtlKey, rtlValue...)
	pipe.Expire(ctx, rtlKey, rtd)

	// allow for maximum 3 tokens simulateously.
	pipe.LTrim(ctx, atlKey, 0, 2)
	pipe.LTrim(ctx, rtlKey, 0, 2)

	_, err := pipe.Exec(ctx)
	return err
}

// LogoutUser deletes access and refresh tokens from Cache and revoke access and refresh tokens by setting them in Cache.
func (cache *RedisStore) LogoutUser(ctx context.Context, accessToken token.Payload, refreshToken token.Payload) error {
	pipe := cache.Client.TxPipeline()

	keys := []string{
		fmt.Sprintf("at:%s:%s", accessToken.UserID.String(), accessToken.ID.String()),
		fmt.Sprintf("rt:%s:%s", refreshToken.UserID.String(), refreshToken.ID.String()),
	}
	pipe.Del(ctx, keys...)

	// Revoke access and refresh token
	revokedAccessTokenKey := fmt.Sprintf("rev:%s:%s", accessToken.UserID.String(), accessToken.ID.String())
	revokedRefreshTokenKey := fmt.Sprintf("rev:%s:%s", refreshToken.UserID.String(), refreshToken.ID.String())
	accessExpiry := accessToken.ExpiredAt.Sub(time.Now().UTC()) + time.Minute
	refreshExpiry := refreshToken.ExpiredAt.Sub(time.Now().UTC()) + time.Minute
	pipe.SetEX(ctx, revokedAccessTokenKey, 1, accessExpiry)
	pipe.SetEX(ctx, revokedRefreshTokenKey, 1, refreshExpiry)

	_, err := pipe.Exec(ctx)
	return err
}

// IsRevoked checks whether a token is revoked by checking the according key `rev:{{userID}}:{{tokenID}}`.
// If the key present, the token is revoked, else perhaps a server error and finally if none of the
// previous then token is not revoked.
func (cache *RedisStore) IsRevoked(ctx context.Context, token token.Payload) (bool, error) {
	key := fmt.Sprintf("rev:%s:%s", token.UserID.String(), token.ID.String())
	_, err := cache.Client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil && err != redis.Nil {
		return true, err
	} else {
		return true, nil
	}
}
