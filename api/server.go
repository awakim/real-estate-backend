package api

import (
	"fmt"

	cache "github.com/awakim/immoblock-backend/cache/redis"
	db "github.com/awakim/immoblock-backend/db/sqlc"
	"github.com/awakim/immoblock-backend/token"
	"github.com/awakim/immoblock-backend/util"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	Config     util.Config
	Store      db.Store
	Cache      cache.Cache
	TokenMaker token.Maker
	Router     *gin.Engine
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config util.Config, store db.Store, cache cache.Cache) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		Config:     config,
		Store:      store,
		Cache:      cache,
		TokenMaker: tokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {

	router := gin.Default()
	router.Use(CORS(server.Config.CorsOrigins))

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/users/refresh", server.refresh)

	authRoutes := router.Group("/").Use(authMiddleware(server.TokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)

	authRoutes.POST("/transfers", server.createTransfer)

	server.Router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}
