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

func (h *Handler) FindAll(c *gin.Context) {
	userID := c.GetString("user_id")

	result, err := h.service.FindAllByUserID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "get attempts success", result)
}

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
