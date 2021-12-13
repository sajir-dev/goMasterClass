package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/sajir-dev/goMasterClass/db/sqlc"
	"github.com/sajir-dev/goMasterClass/utils"
	"github.com/stretchr/testify/assert"
)

func NewTestServer(t *testing.T, store db.Store) *Server {
	config := utils.Config{
		SymmetricKey:  utils.RandomString(32),
		TokenDuration: time.Minute,
	}

	server, err := NewServer(store, config)
	assert.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
