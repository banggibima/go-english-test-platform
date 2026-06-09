package attempt_answers

import (
	"github.com/banggibima/go-english-test-platform/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, handler *Handler, jwtSecret string) {
	answers := router.Group("/attempt-answers")
	answers.Use(middleware.Auth(jwtSecret))

	answers.POST("", handler.SaveOrUpdate)
	answers.GET("", handler.FindByAttemptID)
}
