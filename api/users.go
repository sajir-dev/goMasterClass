package api

import (
	"database/sql"
	"net/http"
	"time"

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

type UserResponse struct {
	Username         string `json:"username"`
	Fullname         string `json:"full_name"`
	Email            string `json:"email"`
	CreatedAt        string `json:"created_at"`
	PasswordChangeAt string `json:"password_changed_at"`
}

func createUserResponse(res db.User) UserResponse {
	return UserResponse{
		Username:         res.Username,
		Email:            res.Email,
		Fullname:         res.FullName,
		CreatedAt:        res.CreatedAt.String(),
		PasswordChangeAt: res.PasswordChangedAt.String(),
	}
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

	resp := createUserResponse(res)

	c.JSON(http.StatusOK, resp)
}

type LoginUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginUserResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

func (s *Server) LoginUser(c *gin.Context) {
	var req LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := s.store.GetUser(c, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err := utils.CheckPassword(req.Password, user.HashedPassword); err != nil {
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, err := s.tokenMaker.CreateToken(req.Username, s.config.TokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, err := s.tokenMaker.CreateRefreshToken(req.Username, time.Hour*24*7)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := LoginUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         createUserResponse(user),
	}

	c.JSON(http.StatusOK, resp)
}
