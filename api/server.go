package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/sajir-dev/goMasterClass/db/sqlc"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "Pong")
	})
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.PUT("/accounts", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	server.router = router
	return server
}

func (s *Server) Start(port string) error {
	return s.router.Run(port)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
