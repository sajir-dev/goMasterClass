package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/sajir-dev/goMasterClass/db/sqlc"
	"github.com/sajir-dev/goMasterClass/utils"
)

// CreateUserParams ...
type CreateUserParams struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type CreateUserResponse struct {
	Username         string `json:"username"`
	Fullname         string `json:"full_name"`
	Email            string `json:"email"`
	CreatedAt        string `json:"created_at"`
	PasswordChangeAt string `json:"password_changed_at"`
}

// CreateUser ...
func (s *Server) CreateUser(c *gin.Context) {
	var req CreateUserParams
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	args := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	res, err := s.store.CreateUser(c, args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				c.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := CreateUserResponse{
		Username:         res.Username,
		Email:            res.Email,
		Fullname:         res.FullName,
		CreatedAt:        res.CreatedAt.String(),
		PasswordChangeAt: res.PasswordChangedAt.String(),
	}

	c.JSON(http.StatusOK, resp)
}
