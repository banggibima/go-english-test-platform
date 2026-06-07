package results

import (
	"net/http"

	"github.com/banggibima/go-english-test-platform/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) FindMyResults(c *gin.Context) {
	userID := c.GetString("user_id")

	result, err := h.service.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "get results success", result)
}

func (h *Handler) FindByID(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	result, err := h.service.FindByID(c.Request.Context(), id, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "get result success", result)
}

func (h *Handler) FindByAttemptID(c *gin.Context) {
	userID := c.GetString("user_id")
	attemptID := c.Query("attempt_id")

	result, err := h.service.FindByAttemptID(c.Request.Context(), attemptID, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "get result by attempt success", result)
}
