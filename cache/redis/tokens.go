package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/awakim/immoblock-backend/token"
)

// DeleteRefreshtoken deletes old refresh tokens
func (cache *RedisStore) DeleteRefreshToken(ctx context.Context, userID string, tokenID string) error {
	key := fmt.Sprintf("rt:%s:%s", userID, tokenID)
	response := cache.Client.Del(ctx, key)
	if err := response.Err(); err != nil {
		return err
	}

	if response.Val() < 1 {
		return fmt.Errorf("invalid Refresh token: %s:%s does not exist", userID, tokenID)
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
	pipe.SetEX(ctx, at, 0, atd)

	if refreshToken.UserID.String() == "" || refreshToken.ID.String() == "" {
		return errors.New("invalid refresh token data to set in cache")
	}
	rt := fmt.Sprintf("rt:%s:%s", refreshToken.UserID.String(), refreshToken.ID.String())
	pipe.SetEX(ctx, rt, 0, rtd)

	atlKey := fmt.Sprintf("atl:%s", accessToken.UserID.String())
	atlValue := []interface{}{accessToken.ID.String()}
	pipe.LPush(ctx, atlKey, atlValue...)

	rtlKey := fmt.Sprintf("rtl:%s", refreshToken.UserID.String())
	rtlValue := []interface{}{refreshToken.ID.String()}
	pipe.LPush(ctx, rtlKey, rtlValue...)

	// allow for maximum 3 tokens simulateously.
	pipe.LTrim(ctx, atlKey, 0, 2)
	pipe.LTrim(ctx, rtlKey, 0, 2)

	_, err := pipe.Exec(ctx)
	return err
}

func (cache *RedisStore) GetTokenData(ctx context.Context, refreshToken token.Payload) error {
	keys := []string{
		fmt.Sprintf("rt:%s:%s", refreshToken.UserID.String(), refreshToken.ID.String()),
		fmt.Sprintf("rtl:%s", refreshToken.UserID.String()),
	}
	resp, err := cache.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return err
	}
	fmt.Println(resp...)
	return nil
}
