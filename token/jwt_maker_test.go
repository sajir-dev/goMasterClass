package token

import (
	"testing"
	"time"

	"github.com/sajir-dev/goMasterClass/utils"
	"github.com/stretchr/testify/assert"
)

func TestJWTCreateToken(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(32))
	assert.Nil(t, err)
	username := utils.RandomOwner()
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	assert.Nil(t, err)
	assert.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	assert.WithinDuration(t, expiredAt, payload.ExpireAt, time.Second)
	assert.Equal(t, username, payload.Username)
	assert.NotEmpty(t, payload.ID)
}

func TestExpiredToken(t *testing.T) {
	maker, err := NewJWTMaker(utils.RandomString(32))
	assert.Nil(t, err)

	username := utils.RandomOwner()
	duration := -time.Minute

	token, err := maker.CreateToken(username, duration)
	assert.Nil(t, err)
	assert.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	assert.NotNil(t, err)
	assert.Nil(t, payload)
}
