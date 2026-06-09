package sections

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
// @Summary Get all sections
// @Description Get list of sections
// @Tags Sections
// @Produce json
// @Success 200 {object} response.Response
// @Router /sections [get]
func (h *Handler) FindAll(c *gin.Context) {
	result, err := h.service.FindAll(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "get sections success", result)
}

// FindByTestID godoc
//
// @Summary Get sections by test ID
// @Description Get sections by test ID
// @Tags Sections
// @Produce json
// @Param test_id query string true "Test ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /sections/by-test [get]
func (h *Handler) FindByTestID(c *gin.Context) {
	testID := c.Query("test_id")

	result, err := h.service.FindByTestID(c.Request.Context(), testID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "get sections by test success", result)
}

// FindByID godoc
//
// @Summary Get section by ID
// @Description Get section detail by ID
// @Tags Sections
// @Produce json
// @Param id path string true "Section ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /sections/{id} [get]
func (h *Handler) FindByID(c *gin.Context) {
	id := c.Param("id")

	result, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "get section success", result)
}

// Create godoc
//
// @Summary Create section
// @Description Create a new section
// @Tags Sections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateSectionRequest true "Create section request"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /sections [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateSectionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if err := validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, "validation failed", err.Error())
		return
	}

	result, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "create section success", result)
}

// Update godoc
//
// @Summary Update section
// @Description Update existing section
// @Tags Sections
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Section ID"
// @Param request body UpdateSectionRequest true "Update section request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /sections/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateSectionRequest

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

	response.Success(c, http.StatusOK, "update section success", result)
}

// Delete godoc
//
// @Summary Delete section
// @Description Delete section by ID
// @Tags Sections
// @Security BearerAuth
// @Produce json
// @Param id path string true "Section ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /sections/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "delete section success", nil)
}
