package attempts

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
	return &Handler{service: service}
}

// StartAttempt godoc
//
// @Summary Start test attempt
// @Description Create a new test attempt
// @Tags Attempts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body StartAttemptRequest true "Start attempt request"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /attempts [post]
func (h *Handler) StartAttempt(c *gin.Context) {
	userID := c.GetString("user_id")

	var req StartAttemptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if err := validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, "validation failed", err.Error())
		return
	}

	result, err := h.service.StartAttempt(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "start attempt success", result)
}

// FindAll godoc
//
// @Summary Get my attempts
// @Description Get all attempts for current user
// @Tags Attempts
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /attempts [get]
func (h *Handler) FindAll(c *gin.Context) {
	userID := c.GetString("user_id")

	result, err := h.service.FindAllByUserID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "get attempts success", result)
}

// FindByID godoc
//
// @Summary Get attempt by ID
// @Description Get attempt detail by ID
// @Tags Attempts
// @Security BearerAuth
// @Produce json
// @Param id path string true "Attempt ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /attempts/{id} [get]
func (h *Handler) FindByID(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	result, err := h.service.FindByID(c.Request.Context(), id, userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "get attempt success", result)
}

// Submit godoc
//
// @Summary Submit attempt
// @Description Submit attempt and publish scoring job to RabbitMQ
// @Tags Attempts
// @Security BearerAuth
// @Produce json
// @Param id path string true "Attempt ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /attempts/{id}/submit [post]
func (h *Handler) Submit(c *gin.Context) {
	userID := c.GetString("user_id")
	id := c.Param("id")

	result, err := h.service.Submit(c.Request.Context(), id, userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "submit attempt success", result)
}
