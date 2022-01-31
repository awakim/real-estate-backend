package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
	}

	newAT, newATID, newRT, newRTID, err := server.TokenMaker.CreateTokenPair(
		refreshToken.UserID,
		server.Config.AccessTokenDuration,
		server.Config.RefreshTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.Cache.SetRefreshToken(ctx, refreshToken.UserID.String(), newATID, server.Config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = server.Cache.SetRefreshToken(ctx, refreshToken.UserID.String(), newRTID, server.Config.RefreshTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := refreshResponse{
		AccessToken:  newAT,
		RefreshToken: newRT,
	}
	ctx.JSON(http.StatusOK, rsp)
}
