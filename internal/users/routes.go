package users

import (
	"github.com/banggibima/go-english-test-platform/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, handler *Handler, jwtSecret string) {
	users := router.Group("/users")
	users.Use(middleware.Auth(jwtSecret))

	users.GET("/me", handler.Me)
	users.PUT("/me", handler.UpdateProfile)
}
