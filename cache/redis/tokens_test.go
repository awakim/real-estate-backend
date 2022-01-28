package cache

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestSetRefreshToken(t *testing.T) {
	uid, err := uuid.NewRandom()
	s := uid.String()

	noErrorCase := struct {
		username   string
		tokenID    string
		expiration time.Duration
	}{
		username:   "validKey",
		tokenID:    s,
		expiration: 2 * time.Second,
	}

	err = testCache.SetRefreshToken(context.Background(), noErrorCase.username, noErrorCase.tokenID, noErrorCase.expiration)
	require.NoError(t, err)

	errorCase := struct {
		username   string
		tokenID    string
		expiration time.Duration
	}{
		username:   "",
		tokenID:    s,
		expiration: 2 * time.Second,
	}
	err = testCache.SetRefreshToken(context.Background(), errorCase.username, errorCase.tokenID, errorCase.expiration)
	require.Error(t, err)
}

func TestDeleteRefreshToken(t *testing.T) {
	uid, _ := uuid.NewRandom()
	s := uid.String()

	noErrorCase := struct {
		username   string
		tokenID    string
		expiration time.Duration
	}{
		username:   "validKey",
		tokenID:    s,
		expiration: 1 * time.Minute,
	}

	_ = testCache.SetRefreshToken(context.Background(), noErrorCase.username, noErrorCase.tokenID, noErrorCase.expiration)
	err := testCache.DeleteRefreshToken(context.Background(), noErrorCase.username, noErrorCase.tokenID)
	require.NoError(t, err)

	errorCase := struct {
		username   string
		tokenID    string
		expiration time.Duration
	}{
		username:   "invalidKey",
		tokenID:    s,
		expiration: 1 * time.Minute,
	}
	err = testCache.DeleteRefreshToken(context.Background(), errorCase.username, errorCase.tokenID)
	require.Error(t, err)
}
