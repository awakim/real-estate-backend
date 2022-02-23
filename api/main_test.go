package api

import (
	"os"
	"testing"
	"time"

	cache "github.com/awakim/immoblock-backend/cache/redis"
	"github.com/awakim/immoblock-backend/config"
	db "github.com/awakim/immoblock-backend/db/sqlc"
	identity "github.com/awakim/immoblock-backend/identity/auth0"
	"github.com/awakim/immoblock-backend/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store, cache cache.Cache, userManager identity.UserManager) *Server {
	config := config.Config{
		TokenSymmetricKey:   util.RandomString(32),
		AccessTokenDuration: time.Minute,
		CorsOrigins: []string{
			"http://localhost:4200",
			"https://localhost:4200",
			"http://localhost:8080",
			"https://localhost:8080",
		},
	}
	server, err := NewServer(config, store, cache, userManager)
	require.NoError(t, err)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
