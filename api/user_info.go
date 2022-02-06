package api

import (
	"errors"
	"fmt"
	"net/http"

	db "github.com/awakim/immoblock-backend/db/sqlc"
	"github.com/awakim/immoblock-backend/token"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

type createUserInfoRequest struct {
	Firstname   string    `json:"firstname" binding:"required,alpha"`
	Lastname    string    `json:"lastname" binding:"required,alpha"`
	PhoneNumber string    `json:"phone_number" binding:"required,e164"`
	Nationality string    `json:"nationality" binding:"required,alpha"`
	Gender      db.Gender `json:"gender" binding:"required,gender"`
	Address     string    `json:"address" binding:"required,ascii"`
	PostalCode  string    `json:"postal_code" binding:"required,alphanum"`
	City        string    `json:"city" binding:"required,alpha"`
	Country     string    `json:"country" binding:"required,alpha"`
}

type userInfoResponse struct {
	UserID      uuid.UUID `json:"user_id"`
	Firstname   string    `json:"firstname"`
	Lastname    string    `json:"lastname"`
	PhoneNumber string    `json:"phone_number"`
	Nationality string    `json:"nationality"`
	Gender      db.Gender `json:"gender"`
	Address     string    `json:"address"`
	PostalCode  string    `json:"postal_code"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
}

func newUserInfoResponse(userInfo db.UserInformation) userInfoResponse {
	return userInfoResponse{
		UserID:      userInfo.UserID,
		Firstname:   userInfo.Firstname,
		Lastname:    userInfo.Lastname,
		PhoneNumber: userInfo.PhoneNumber,
		Nationality: userInfo.Nationality,
		Gender:      userInfo.Gender,
		Address:     userInfo.Address,
		PostalCode:  userInfo.PostalCode,
		City:        userInfo.City,
		Country:     userInfo.Country,
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

	arg := db.CreateUserInfoParams{
		UserID:      authPayload.UserID,
		Firstname:   req.Firstname,
		Lastname:    req.Lastname,
		PhoneNumber: req.PhoneNumber,
		Nationality: req.Nationality,
		Gender:      req.Gender,
		Address:     req.Address,
		PostalCode:  req.PostalCode,
		City:        req.City,
		Country:     req.Country,
	}

	userInfo, err := server.Store.CreateUserInfo(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				fmt.Println(req)
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
