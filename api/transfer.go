package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/sajir-dev/goMasterClass/db/sqlc"
	"github.com/sajir-dev/goMasterClass/token"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (s *Server) createTransfer(c *gin.Context) {
	var req transferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, exists := s.validAccount(c, req.FromAccountID, req.Currency)
	if !exists {
		return
	}

	_, exists = s.validAccount(c, req.ToAccountID, req.Currency)
	if !exists {
		return
	}

	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account id does not belong to the authenticated user")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	args := db.TransferTxParams{
		From_Account_ID: req.FromAccountID,
		To_Account_ID:   req.ToAccountID,
		Amount:          req.Amount,
	}

	result, err := s.store.TransferTx(c, args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, result)
}

func (s *Server) validAccount(c *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := s.store.GetAccount(c, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return account, false
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return account, false
	}

	if account.Currency != currency {
		c.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("wrong currency:%s input given for the account: %v, actual currency: %s", currency, accountID, account.Currency)))
		return account, false
	}

	return account, true
}
