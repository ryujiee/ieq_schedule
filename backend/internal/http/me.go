package http

import (
	"net/http"

	"ieq/backend/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func meHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get(middlewares.CtxUserIDKey)
		role, _ := c.Get(middlewares.CtxRoleKey)

		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
			"role":    role,
		})
	}
}
