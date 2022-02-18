package api

import (
	"database/sql"
	"errors"
	"net/http"

	db "github.com/awakim/immoblock-backend/db/sqlc"
	"github.com/awakim/immoblock-backend/token"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

type createUserInfoRequest struct {
	Firstname   string `json:"firstname" binding:"required,alpha"`
	Lastname    string `json:"lastname" binding:"required,alpha"`
	PhoneNumber string `json:"phone_number" binding:"required,e164"`
	Nationality string `json:"nationality" binding:"required,alpha"`
	Address     string `json:"address" binding:"required,ascii"`
	PostalCode  string `json:"postal_code" binding:"required,alphanum"`
	City        string `json:"city" binding:"required,alpha"`
	Country     string `json:"country" binding:"required,alpha"`
}

type userInfoResponse struct {
	UserID           uuid.UUID `json:"user_id"`
	Firstname        string    `json:"firstname"`
	Lastname         string    `json:"lastname"`
	PhoneNumber      string    `json:"phone_number"`
	Nationality      string    `json:"nationality"`
	Address          string    `json:"address"`
	PostalCode       string    `json:"postal_code"`
	City             string    `json:"city"`
	Country          string    `json:"country"`
	VerificationStep int16     `json:"verification_step"`
}

func newUserInfoResponse(userInfo db.UserInformation) userInfoResponse {
	return userInfoResponse{
		UserID:           userInfo.UserID,
		Firstname:        userInfo.Firstname,
		Lastname:         userInfo.Lastname,
		PhoneNumber:      userInfo.PhoneNumber,
		Nationality:      userInfo.Nationality,
		Address:          userInfo.Address,
		PostalCode:       userInfo.PostalCode,
		City:             userInfo.City,
		Country:          userInfo.Country,
		VerificationStep: userInfo.VerificationStep,
	}
}

func (server *Server) createUserInfo(ctx *gin.Context) {
	var req createUserInfoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			ctx.JSON(http.StatusBadRequest, gin.H{"errors": ValidationError(verr)})
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	exists, err := server.Store.ExistsUserInfo(ctx, authPayload.UserID)
	if err == nil && exists {
		errRowAlreadyExist := errors.New("user information already provided please contact support")
		ctx.JSON(http.StatusForbidden, errorResponse(errRowAlreadyExist))
		return
	} else if err != nil && err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserInfoParams{
		UserID:           authPayload.UserID,
		Firstname:        req.Firstname,
		Lastname:         req.Lastname,
		PhoneNumber:      req.PhoneNumber,
		Nationality:      req.Nationality,
		Address:          req.Address,
		PostalCode:       req.PostalCode,
		City:             req.City,
		Country:          req.Country,
		VerificationStep: 1,
	}

	userInfo, err := server.Store.CreateUserInfo(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				errPhoneAlreadyExists := errors.New("this phone number already exists")
				ctx.JSON(http.StatusForbidden, errorResponse(errPhoneAlreadyExists))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserInfoResponse(userInfo)
	ctx.JSON(http.StatusOK, rsp)
}

func (server *Server) getUserInfo(ctx *gin.Context) {

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	userInfo, err := server.Store.GetUserInfo(ctx, authPayload.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(errors.New("user has not provided information yet")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := newUserInfoResponse(userInfo)
	ctx.JSON(http.StatusOK, rsp)
}
