package attempts

import (
	"github.com/banggibima/go-english-test-platform/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, handler *Handler, jwtSecret string) {
	attempts := router.Group("/attempts")
	attempts.Use(middleware.Auth(jwtSecret))

	attempts.POST("", handler.StartAttempt)
	attempts.GET("", handler.FindAll)
	attempts.GET("/:id", handler.FindByID)
	attempts.POST("/:id/submit", handler.Submit)
}
