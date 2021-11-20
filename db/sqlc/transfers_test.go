package db

import (
	"context"
	"testing"

	"github.com/sajir-dev/goMasterClass/utils"
	"github.com/stretchr/testify/assert"
)

func Test_CreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	args := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        utils.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	assert.Nil(t, err)
	assert.NotZero(t, transfer.ID)
	assert.Equal(t, transfer.Amount, args.Amount)
	assert.Equal(t, transfer.FromAccountID, account1.ID)
	assert.Equal(t, transfer.ToAccountID, account2.ID)
}
