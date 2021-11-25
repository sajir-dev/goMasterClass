package db

import (
	"context"
	"testing"

	"github.com/sajir-dev/goMasterClass/utils"
	"github.com/stretchr/testify/assert"
)

func Test_CreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func CreateRandomUser(t *testing.T) User {
	args := CreateUserParams{
		Username:       utils.RandomString(6),
		HashedPassword: "somepassword",
		FullName:       utils.RandomOwner(),
		Email:          utils.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), args)
	assert.NoError(t, err)
	assert.Equal(t, user.Username, args.Username)
	assert.Equal(t, user.Email, args.Email)
	assert.Equal(t, user.HashedPassword, args.HashedPassword)
	assert.Equal(t, user.FullName, args.FullName)
	assert.NotEmpty(t, user.CreatedAt)
	assert.True(t, user.PasswordChangedAt.IsZero())
	return user
}

func Test_GetUser(t *testing.T) {
	user1 := CreateRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	assert.NoError(t, err)
	assert.Equal(t, user1, user2)
}
