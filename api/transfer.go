package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/awakim/immoblock-backend/db/sqlc"
	"github.com/awakim/immoblock-backend/token"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type transferRequest struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,min=1"`
	Amount        int64 `json:"amount" binding:"required,gt=0"`
	PropertyID    int64 `json:"property_id" binding:"required,min=1"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			ctx.JSON(http.StatusBadRequest, gin.H{"errors": ValidationError(verr)})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"errors": errorResponse(err)})
		return
	}

	fromAccount, valid := server.validAccount(ctx, req.FromAccountID, req.PropertyID)
	if !valid {
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.UserID != authPayload.UserID {
		err := errors.New("from account does not belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	_, valid = server.validAccount(ctx, req.ToAccountID, req.PropertyID)
	if !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.Store.TransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64, propertyID int64) (db.Account, bool) {
	account, err := server.Store.GetAccount(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	_, err = server.Store.GetProperty(ctx, propertyID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.PropertyID != propertyID {
		err := fmt.Errorf("account [%d] property_id mismatch: %v vs %v", accountID, account.PropertyID, propertyID)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}

	return account, true
}
