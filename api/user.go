package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/auth0/go-auth0/management"
	db "github.com/awakim/immoblock-backend/db/sqlc"
	"github.com/awakim/immoblock-backend/util"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

type createUserRequest struct {
	Nickname string `json:"nickname" binding:"required,alphanum,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type userResponse struct {
	UserID            uuid.UUID `json:"user_id"`
	Nickname          string    `json:"nickname"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		UserID:            user.ID,
		Nickname:          user.Nickname,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			ctx.JSON(http.StatusBadRequest, gin.H{"errors": ValidationError(verr)})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errorResponse(err)})
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Email:          req.Email,
		Nickname:       req.Nickname,
		HashedPassword: hashedPassword,
	}

	user, err := server.Store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				errEmailAlreadyExists := errors.New("this email already exists")
				ctx.JSON(http.StatusForbidden, errorResponse(errEmailAlreadyExists))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	uid := user.ID
	strUUID := uid.String()
	verifyEmail := true
	connection := "Username-Password-Authentication"
	authZeroUser := &management.User{
		ID:          &strUUID,
		Name:        &req.Nickname,
		Email:       &req.Email,
		Password:    &hashedPassword,
		VerifyEmail: &verifyEmail,
		Connection:  &connection,
	}
	err = server.UserManager.Create(authZeroUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginUserResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			ctx.JSON(http.StatusBadRequest, gin.H{"errors": ValidationError(verr)})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errorResponse(err)})
		return
	}

	user, err := server.Store.GetUser(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid credentials")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid credentials")))
		return
	}

	newAT, newATST, newRT, newRTST, err := server.TokenMaker.CreateTokenPair(
		user.ID,
		user.IsAdmin,
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

	rsp := loginUserResponse{
		AccessToken:  newATST,
		RefreshToken: newRTST,
		User:         newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, rsp)
}

type logoutUserRequest struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (server *Server) logoutUser(ctx *gin.Context) {
	var req logoutUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			ctx.JSON(http.StatusBadRequest, gin.H{"errors": ValidationError(verr)})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errorResponse(err)})
		return
	}

	accessToken, err := server.TokenMaker.VerifyToken(req.AccessToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	refreshToken, err := server.TokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	err = server.Cache.LogoutUser(ctx, *accessToken, *refreshToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user has successfully logged out",
	})
}
