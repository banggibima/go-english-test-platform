package dashboard

import (
	"github.com/banggibima/go-english-test-platform/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, handler *Handler, jwtSecret string) {
	dashboard := router.Group("/dashboard")
	dashboard.Use(middleware.Auth(jwtSecret))

	dashboard.GET("", handler.GetSummary)
}
