package db

import (
	"context"
	"testing"

	"github.com/sajir-dev/goMasterClass/utils"
	"github.com/stretchr/testify/assert"
)

func Test_CreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	args := CreateEntryParams{
		AccountID: account.ID,
		Amount:    utils.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), args)
	assert.Nil(t, err)
	assert.Equal(t, entry.AccountID, account.ID)
	assert.Equal(t, entry.Amount, args.Amount)
	assert.NotEmpty(t, entry.ID)
}

func Test_DeleteEntry(t *testing.T) {
	account := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	assert.Nil(t, err)
}

func Test_ListEntries(t *testing.T) {
	account := createRandomAccount(t)
	amounts := []int64{
		utils.RandomMoney(),
		utils.RandomMoney(),
		utils.RandomMoney(),
		utils.RandomMoney(),
		utils.RandomMoney(),
		utils.RandomMoney(),
		utils.RandomMoney(),
		utils.RandomMoney(),
		utils.RandomMoney(),
		utils.RandomMoney(),
	}
	for i := 0; i < 10; i++ {
		testQueries.CreateEntry(context.Background(), CreateEntryParams{
			AccountID: account.ID,
			Amount:    amounts[i],
		})
	}

	entries, err := testQueries.ListEntries(context.Background(), ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	})
	assert.Nil(t, err)
	assert.Len(t, entries, 5)
	for k, v := range entries {
		assert.NotEmpty(t, v.ID)
		assert.Equal(t, v.AccountID, account.ID)
		assert.Equal(t, v.Amount, amounts[k+5])
	}
}

func Test_UpdateEntry(t *testing.T) {
	account1 := createRandomAccount(t)
	balance := utils.RandomMoney()
	account2, err := testQueries.UpdateAccount(context.Background(), UpdateAccountParams{
		ID:      account1.ID,
		Balance: balance,
	})

	assert.Nil(t, err)
	assert.Equal(t, account1.ID, account2.ID)
	assert.Equal(t, account2.Balance, balance)
}
