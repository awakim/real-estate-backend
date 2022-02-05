package cache

import (
	"context"
	"testing"
	"time"

	"github.com/awakim/immoblock-backend/token"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSetTokenData(t *testing.T) {
	uid1, _ := uuid.NewRandom()
	uid2, _ := uuid.NewRandom()

	at, _ := token.NewPayload(uid1, false, time.Minute)
	rt, _ := token.NewPayload(uid2, false, time.Minute)

	noErrorCase := struct {
		accessToken  *token.Payload
		refreshToken *token.Payload
	}{
		accessToken:  at,
		refreshToken: rt,
	}

	err := testCache.SetTokenData(context.Background(), *noErrorCase.accessToken, time.Minute, *noErrorCase.refreshToken, time.Minute)
	require.NoError(t, err)
}

func TestDeleteRefreshToken(t *testing.T) {
	uid1, _ := uuid.NewRandom()
	uid2, _ := uuid.NewRandom()

	at, _ := token.NewPayload(uid1, false, time.Minute)
	rt, _ := token.NewPayload(uid2, false, time.Minute)

	noErrorCase := struct {
		accessToken  *token.Payload
		refreshToken *token.Payload
	}{
		accessToken:  at,
		refreshToken: rt,
	}

	_ = testCache.SetTokenData(context.Background(), *noErrorCase.accessToken, time.Minute, *noErrorCase.refreshToken, time.Minute)
	err := testCache.DeleteRefreshToken(context.Background(), rt.UserID.String(), rt.ID.String())
	require.NoError(t, err)

	err = testCache.DeleteRefreshToken(context.Background(), "invalid_key", rt.ID.String())
	require.Error(t, err)
}
