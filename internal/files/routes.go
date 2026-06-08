package files

import (
	"github.com/banggibima/go-english-test-platform/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup, handler *Handler, jwtSecret string) {
	files := router.Group("/files")
	files.Use(middleware.Auth(jwtSecret))

	files.POST("/upload", handler.Upload)
	files.GET("", handler.FindMyFiles)
	files.GET("/:id", handler.FindByID)
}
