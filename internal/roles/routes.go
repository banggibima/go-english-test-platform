package roles

import (
	"github.com/banggibima/go-english-test-platform/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, handler *Handler, jwtSecret string) {
	roles := router.Group("/roles")
	roles.Use(middleware.Auth(jwtSecret))
	roles.Use(middleware.Role("admin"))

	roles.GET("", handler.FindAll)
	roles.GET("/:id", handler.FindByID)
}
