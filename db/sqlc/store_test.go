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

	account1, _ = testQueries.UpdateAccount(context.Background(), UpdateAccountParams{
		ID:      account1.ID,
		Balance: int64(500),
	})

	// run n concurrent transfer txns
	n := 5
	amount := int64(10)
	existed := make(map[int64]bool)

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
		fromAccount := res.FromAccount
		assert.NotEmpty(t, fromAccount)
		assert.Equal(t, account1.ID, fromAccount.ID)

		toAccount := res.ToAccount
		assert.NotEmpty(t, toAccount)
		assert.Equal(t, toAccount.ID, account2.ID)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		assert.Equal(t, diff1, diff2)
		assert.True(t, diff1 > 0)
		assert.True(t, diff1%amount == 0)

		k := diff1 / amount
		assert.True(t, k >= 1 && k <= int64(n))
		assert.NotContains(t, existed, k)
		existed[k] = true
	}

	// check the final updated balance
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.Nil(t, err)
	assert.NotEmpty(t, updateAccount1)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	assert.Nil(t, err)
	assert.NotEmpty(t, updateAccount2)

	assert.Equal(t, updateAccount1.Balance, account1.Balance-int64(n)*amount)
	assert.Equal(t, updateAccount2.Balance, account2.Balance+int64(n)*amount)
}

func TestTransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	account1, _ = testQueries.UpdateAccount(context.Background(), UpdateAccountParams{
		ID:      account1.ID,
		Balance: int64(500),
	})

	// run n concurrent transfer txns
	n := 10
	amount := int64(10)

	errs := make(chan error)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			args := TransferTxParams{
				From_Account_ID: account1.ID,
				To_Account_ID:   account2.ID,
				Amount:          amount,
			}
			if i%2 == 1 {
				args = TransferTxParams{
					From_Account_ID: account2.ID,
					To_Account_ID:   account1.ID,
					Amount:          amount,
				}
			}
			_, err := store.TransferTx(context.Background(), args)
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs

		assert.Nil(t, err)
	}

	// check the final updated balance
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.Nil(t, err)
	assert.NotEmpty(t, updateAccount1)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	assert.Nil(t, err)
	assert.NotEmpty(t, updateAccount2)

	assert.Equal(t, updateAccount1.Balance, account1.Balance)
	assert.Equal(t, updateAccount2.Balance, account2.Balance)
}
