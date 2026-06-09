package dashboard

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

// GetSummary godoc
//
// @Summary Get dashboard summary
// @Description Get dashboard summary for current user
// @Tags Dashboard
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /dashboard [get]
func (h *Handler) GetSummary(c *gin.Context) {
	userID := c.GetString("user_id")

	result, err := h.service.GetSummary(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "get dashboard success", result)
}
