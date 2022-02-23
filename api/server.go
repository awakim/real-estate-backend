package api

import (
	"fmt"
	"reflect"
	"strings"

	cache "github.com/awakim/immoblock-backend/cache/redis"
	"github.com/awakim/immoblock-backend/config"
	db "github.com/awakim/immoblock-backend/db/sqlc"
	identity "github.com/awakim/immoblock-backend/identity/auth0"
	"github.com/awakim/immoblock-backend/token"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	Config      config.Config
	Store       db.Store
	Cache       cache.Cache
	TokenMaker  token.Maker
	Router      *gin.Engine
	UserManager identity.UserManager
}

// NewServer creates a new HTTP server and set up routing.
func NewServer(config config.Config, store db.Store, cache cache.Cache, userManager identity.UserManager) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		Config:      config,
		Store:       store,
		Cache:       cache,
		TokenMaker:  tokenMaker,
		UserManager: userManager,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {

	router := gin.Default()
	router.Use(CORS(server.Config.CorsOrigins))

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginRateLimiter, server.loginUser)
	router.POST("/users/refresh", server.refresh)

	authRoutes := router.Group("/").Use(auth(server.TokenMaker), server.revoked)

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)

	authRoutes.POST("/transfers", server.createTransfer)

	authRoutes.GET("/users/info", server.getUserInfo)
	authRoutes.POST("/users/info", server.createUserInfo)
	authRoutes.POST("/users/logout", server.logoutUser)

	server.Router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.Router.Run(address)
}
