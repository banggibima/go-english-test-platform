package middleware

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func Role(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentRole := c.GetString("role")

		if slices.Contains(roles, currentRole) {
			c.Next()
			return
		}

		c.JSON(http.StatusForbidden, gin.H{
			"message": "forbidden",
		})
		c.Abort()
	}
}
