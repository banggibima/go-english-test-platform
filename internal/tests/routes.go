package tests

import (
	"github.com/banggibima/go-english-test-platform/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, handler *Handler, jwtSecret string) {
	tests := router.Group("/tests")

	tests.GET("", handler.FindAll)
	tests.GET("/:id", handler.FindByID)

	admin := tests.Group("")
	admin.Use(middleware.Auth(jwtSecret))
	admin.Use(middleware.Role("admin"))

	admin.POST("", handler.Create)
	admin.PUT("/:id", handler.Update)
	admin.DELETE("/:id", handler.Delete)
}
