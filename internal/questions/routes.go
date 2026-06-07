package questions

import (
	"github.com/banggibima/go-english-test-platform/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, handler *Handler, jwtSecret string) {
	questions := router.Group("/questions")

	questions.GET("", handler.FindAll)
	questions.GET("/by-section", handler.FindBySectionID)
	questions.GET("/:id", handler.FindByID)

	admin := questions.Group("")
	admin.Use(middleware.Auth(jwtSecret))
	admin.Use(middleware.Role("admin"))

	admin.POST("", handler.Create)
	admin.PUT("/:id", handler.Update)
	admin.DELETE("/:id", handler.Delete)
}
