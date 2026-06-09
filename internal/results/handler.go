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

// FindMyResults godoc
//
// @Summary Get my results
// @Description Get all results for current user
// @Tags Results
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /results [get]
func (h *Handler) FindMyResults(c *gin.Context) {
	userID := c.GetString("user_id")

	result, err := h.service.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "get results success", result)
}

// FindByAttemptID godoc
//
// @Summary Get result by attempt ID
// @Description Get result detail by attempt ID
// @Tags Results
// @Security BearerAuth
// @Produce json
// @Param attempt_id query string true "Attempt ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /results/by-attempt [get]
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

// FindByID godoc
//
// @Summary Get result by ID
// @Description Get result detail by ID
// @Tags Results
// @Security BearerAuth
// @Produce json
// @Param id path string true "Result ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /results/{id} [get]
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
