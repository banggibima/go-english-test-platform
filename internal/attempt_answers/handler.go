package attempt_answers

import (
	"net/http"

	"github.com/banggibima/go-english-test-platform/pkg/response"
	"github.com/banggibima/go-english-test-platform/pkg/validator"
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

// SaveOrUpdate godoc
//
// @Summary Save or update attempt answer
// @Description Save or update answer for an attempt question
// @Tags Attempt Answers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body SaveAnswerRequest true "Save attempt answer request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /attempt-answers [post]
func (h *Handler) SaveOrUpdate(c *gin.Context) {
	var req SaveAnswerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if err := validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, "validation failed", err.Error())
		return
	}

	result, err := h.service.SaveOrUpdate(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "save answer success", result)
}

// FindByAttemptID godoc
//
// @Summary Get answers by attempt ID
// @Description Get all answers for an attempt
// @Tags Attempt Answers
// @Security BearerAuth
// @Produce json
// @Param attempt_id query string true "Attempt ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /attempt-answers [get]
func (h *Handler) FindByAttemptID(c *gin.Context) {
	attemptID := c.Query("attempt_id")

	result, err := h.service.FindByAttemptID(c.Request.Context(), attemptID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "get attempt answers success", result)
}
