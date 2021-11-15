package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// run n concurrent transfer txns
	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		go func() {
			args := TransferTxParams{
				From_Account_ID: account1.ID,
				To_Account_ID:   account2.ID,
				Amount:          amount,
			}
			res, err := store.TransferTx(context.Background(), args)
			errs <- err
			results <- res
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		res := <-results

		assert.Nil(t, err)
		assert.NotNil(t, res)

		// check transfer
		assert.NotEmpty(t, res.Transfer)
		assert.Equal(t, account1.ID, res.Transfer.FromAccountID)
		assert.Equal(t, account2.ID, res.Transfer.ToAccountID)
		assert.Equal(t, res.Transfer.Amount, amount)
		assert.NotZero(t, res.Transfer.ID)
		assert.NotZero(t, res.Transfer.CreatedAt)

		// check entries
		assert.Equal(t, account1.ID, res.FromEntry.AccountID)
		assert.Equal(t, account2.ID, res.ToEntry.AccountID)
		assert.Equal(t, -amount, res.FromEntry.Amount)
		assert.Equal(t, amount, res.ToEntry.Amount)
		assert.NotEmpty(t, res.FromEntry.CreatedAt)
		assert.NotEmpty(t, res.ToEntry.CreatedAt)
		assert.NotEmpty(t, res.FromEntry.ID)
		assert.NotEmpty(t, res.ToEntry.ID)

		_, err = store.GetEntry(context.Background(), res.FromEntry.ID)
		assert.Nil(t, err)

		_, err = store.GetEntry(context.Background(), res.ToEntry.ID)
		assert.Nil(t, err)

		// TODO: check account balance
	}
}
