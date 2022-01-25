package api

import (
	"os"
	"testing"
	"time"

	cache "github.com/awakim/immoblock-backend/cache/redis"
	db "github.com/awakim/immoblock-backend/db/sqlc"
	"github.com/awakim/immoblock-backend/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store, cache cache.Cache) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}
	server, err := NewServer(config, store, cache)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
