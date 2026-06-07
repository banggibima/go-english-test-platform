package results

import (
	"github.com/banggibima/go-english-test-platform/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, handler *Handler, jwtSecret string) {
	results := router.Group("/results")
	results.Use(middleware.Auth(jwtSecret))

	results.GET("", handler.FindMyResults)
	results.GET("/by-attempt", handler.FindByAttemptID)
	results.GET("/:id", handler.FindByID)
}
