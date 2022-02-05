package token

import (
	"testing"
	"time"

	"github.com/awakim/immoblock-backend/util"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPasetoMaker(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	userID, err := uuid.NewRandom()
	require.NoError(t, err)
	isAdmin := false
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, ss, err := maker.CreateToken(userID, isAdmin, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, ss)

	payload, err := maker.VerifyToken(ss)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	require.NotZero(t, payload.ID)
	require.Equal(t, userID, payload.UserID)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)
}

func TestExpiredPasetoToken(t *testing.T) {
	maker, err := NewPasetoMaker(util.RandomString(32))
	require.NoError(t, err)

	uid, err := uuid.NewRandom()
	require.NoError(t, err)
	isAdmin := false
	token, ss, err := maker.CreateToken(uid, isAdmin, -time.Minute)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, ss)

	payload, err := maker.VerifyToken(ss)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredToken.Error())
	require.Nil(t, payload)
}
