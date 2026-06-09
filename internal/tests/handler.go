package tests

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

// FindAll godoc
//
// @Summary Get all tests
// @Description Get list of available tests
// @Tags Tests
// @Produce json
// @Success 200 {object} response.Response
// @Router /tests [get]
func (h *Handler) FindAll(c *gin.Context) {
	result, err := h.service.FindAll(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "get tests success", result)
}

// FindByID godoc
//
// @Summary Get test by ID
// @Description Get test detail by ID
// @Tags Tests
// @Produce json
// @Param id path string true "Test ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /tests/{id} [get]
func (h *Handler) FindByID(c *gin.Context) {
	id := c.Param("id")

	result, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "get test success", result)
}

// Create godoc
//
// @Summary Create test
// @Description Create a new test
// @Tags Tests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateTestRequest true "Create test request"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /tests [post]
func (h *Handler) Create(c *gin.Context) {
	userID := c.GetString("user_id")

	var req CreateTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if err := validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, "validation failed", err.Error())
		return
	}

	result, err := h.service.Create(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "create test success", result)
}

// Update godoc
//
// @Summary Update test
// @Description Update existing test
// @Tags Tests
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Test ID"
// @Param request body UpdateTestRequest true "Update test request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /tests/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if err := validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, "validation failed", err.Error())
		return
	}

	result, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "update test success", result)
}

// Delete godoc
//
// @Summary Delete test
// @Description Delete test by ID
// @Tags Tests
// @Security BearerAuth
// @Produce json
// @Param id path string true "Test ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /tests/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "delete test success", nil)
}
