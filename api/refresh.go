package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type refreshRequest struct {
	RefreshTokenString string `json:"refresh_token" binding:"required"`
}

type refreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (server *Server) refresh(ctx *gin.Context) {
	var req refreshRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	refreshToken, err := server.TokenMaker.VerifyToken(req.RefreshTokenString)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	prevTokenID := refreshToken.ID.String()
	if prevTokenID != "" {
		if err := server.Cache.DeleteRefreshToken(ctx, refreshToken.UserID.String(), prevTokenID); err != nil {
			if err == redis.Nil {
				ctx.JSON(http.StatusNotFound, errorResponse(errors.New("refresh not found")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	newAT, newATST, newRT, newRTST, err := server.TokenMaker.CreateTokenPair(
		refreshToken.UserID,
		refreshToken.IsAdmin,
		server.Config.AccessTokenDuration,
		server.Config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.Cache.SetTokenData(ctx, newAT, server.Config.AccessTokenDuration, newRT, server.Config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := refreshResponse{
		AccessToken:  newATST,
		RefreshToken: newRTST,
	}
	ctx.JSON(http.StatusOK, rsp)
}
