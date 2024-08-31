package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"example.com/simple_bank/db/sqlc"
	"example.com/simple_bank/token"
	"github.com/gin-gonic/gin"
)

type transferParams struct {
	FromAccountId int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,validateCurrency"`
}

func (server *Server) transfer(ctx *gin.Context) {

	var req transferParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := sqlc.TransferTxParams{
		FromAccountId: req.FromAccountId,
		ToAccountId:   req.ToAccountId,
		Amount:        req.Amount,
	}

	fromAccount, isValid := server.isValidAccount(ctx, req.FromAccountId, req.Currency)
	if !isValid {
		return
	}

	autyPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if autyPayload.Username != fromAccount.Owner {
		err := errors.New("fromAccount is not of authenticated user")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	_, isValid = server.isValidAccount(ctx, req.ToAccountId, req.Currency)

	if !isValid {
		return
	}

	res, err := server.store.TransferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (server *Server) isValidAccount(ctx *gin.Context, accountId int64, currency string) (sqlc.Account, bool) {

	account, err := server.store.GetAccount(ctx, accountId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err = fmt.Errorf("account %d has currency mismatch. %s vs %s", accountId, account.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return account, false
	}
	return account, true
}
