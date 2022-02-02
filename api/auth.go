package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/awakim/immoblock-backend/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "Bearer"
	authorizationPayloadKey = "X-Auth-Payload"
)

// authMiddleware creates a gin middleware for authorization
func auth(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := ""
		_, err := fmt.Sscanf(authorizationHeader, "Bearer %s", &accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("invalid Token")))
			return
		}

		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

func (server *Server) revoked(ctx *gin.Context) {

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	b, err := server.Cache.IsRevoked(ctx, *payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	if b {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(errors.New("unauthorized")))
		return
	}
	ctx.Set(authorizationPayloadKey, payload)
	ctx.Next()
}
