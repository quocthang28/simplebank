package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	db "simplebank/db/sqlc"
	"simplebank/token"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) getAccountById(ctx *gin.Context, id int64) db.Account {
	account, err := server.store.GetAccount(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			NotFound(ctx, errorResponse(err))
			return db.Account{}
		}

		InternalServerError(ctx, errorResponse(err))
		return db.Account{}
	}

	return account
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		BadRequest(ctx, errorResponse(err))
		return
	}

	senderAccount := server.getAccountById(ctx, req.FromAccountID)
	receiverAccount := server.getAccountById(ctx, req.ToAccountID)

	if !server.checkSenderUsername(ctx, senderAccount.Owner) {
		return
	}

	if !server.checkSenderAccountBalance(ctx, senderAccount.Balance, req.Amount) {
		return
	}

	if !server.validAccount(ctx, senderAccount.Currency, req.Currency) {
		return
	}

	if !server.validAccount(ctx, receiverAccount.Currency, req.Currency) {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(ctx, arg)
	if err != nil {
		InternalServerError(ctx, errorResponse(err))
		return
	}

	OK(ctx, result)
}

func (server *Server) validAccount(ctx *gin.Context, accountCurrency, currency string) bool {
	if accountCurrency != currency {
		err := fmt.Errorf("account currency mismatch: %s vs %s", accountCurrency, currency)
		BadRequest(ctx, errorResponse(err))
		return false
	}

	return true
}

func (server *Server) checkSenderAccountBalance(ctx *gin.Context, balance, amount int64) bool {
	if balance <= 0 || balance < amount {
		err := errors.New("requested account insufficient balance")
		BadRequest(ctx, errorResponse(err))
		return false
	}

	return true
}

func (server *Server) checkSenderUsername(ctx *gin.Context, owner string) bool {
	authPayload := ctx.MustGet(authPayloadKey).(*token.Payload)
	if owner != authPayload.Username {
		err := errors.New("requested account does not belong to the authenticated user")
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
	}

	return true
}
