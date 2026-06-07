package sections

import (
	"github.com/banggibima/go-english-test-platform/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, handler *Handler, jwtSecret string) {
	sections := router.Group("/sections")

	sections.GET("", handler.FindAll)
	sections.GET("/by-test", handler.FindByTestID)
	sections.GET("/:id", handler.FindByID)

	admin := sections.Group("")
	admin.Use(middleware.Auth(jwtSecret))
	admin.Use(middleware.Role("admin"))

	admin.POST("", handler.Create)
	admin.PUT("/:id", handler.Update)
	admin.DELETE("/:id", handler.Delete)
}
