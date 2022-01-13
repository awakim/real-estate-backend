package api

import (
	db "github.com/awakim/immoblock-backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP requests for the immoblock service.
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new server and setup routing.
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse returns a gin.H containing the error
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
