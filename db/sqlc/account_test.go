package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/sajir-dev/goMasterClass/utils"
	"github.com/stretchr/testify/assert"
)

func createRandomAccount(t *testing.T) Account {
	user := CreateRandomUser(t)
	params := CreateAccountParams{
		Owner:    user.Username,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	res, err := testQueries.CreateAccount(context.Background(), params)
	assert.Equal(t, res.Owner, params.Owner)
	assert.Equal(t, res.Balance, params.Balance)
	assert.Equal(t, res.Currency, params.Currency)
	assert.NotEmpty(t, res.CreatedAt)
	assert.NotEmpty(t, res.ID)
	assert.Nil(t, err)
	return res
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.ID)
	assert.Nil(t, err)
	assert.Equal(t, account1, account2)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	req := UpdateAccountParams{
		ID:      account1.ID,
		Balance: utils.RandomMoney(),
	}

	updatedAccount, err1 := testQueries.UpdateAccount(context.Background(), req)
	account2, err2 := testQueries.GetAccount(context.Background(), req.ID)
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.Equal(t, updatedAccount, account2)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account.ID)
	assert.Nil(t, err)

	account, err = testQueries.GetAccount(context.Background(), account.ID)
	assert.NotNil(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Empty(t, account)
}

func TestListAccount(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createRandomAccount(t)
	}

	params := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), params)
	assert.Nil(t, err)
	assert.NotEmpty(t, accounts)
	for _, account := range accounts {
		assert.NotEmpty(t, account)
		assert.Equal(t, lastAccount.Owner, account.Owner)
	}
}
