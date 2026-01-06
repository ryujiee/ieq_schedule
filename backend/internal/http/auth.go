package http

import (
	"net/http"

	"ieq/backend/internal/services"

	"github.com/gin-gonic/gin"
)

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResp struct {
	Token string `json:"token"`
	User  struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
		Role  string `json:"role"`
	} `json:"user"`
}

func loginHandler(authSvc *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		token, user, err := authSvc.Login(req.Email, req.Password)
		if err != nil {
			if err == services.ErrInvalidCredentials {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		var resp loginResp
		resp.Token = token
		resp.User.ID = user.ID
		resp.User.Email = user.Email
		resp.User.Name = user.Name
		resp.User.Role = user.Role.Slug

		c.JSON(http.StatusOK, resp)
	}
}
