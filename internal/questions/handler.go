package questions

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
// @Summary Get all questions
// @Description Get list of questions
// @Tags Questions
// @Produce json
// @Success 200 {object} response.Response
// @Router /questions [get]
func (h *Handler) FindAll(c *gin.Context) {
	result, err := h.service.FindAll(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "get questions success", result)
}

// FindBySectionID godoc
//
// @Summary Get questions by section ID
// @Description Get questions by section ID
// @Tags Questions
// @Produce json
// @Param section_id query string true "Section ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /questions/by-section [get]
func (h *Handler) FindBySectionID(c *gin.Context) {
	sectionID := c.Query("section_id")

	result, err := h.service.FindBySectionID(c.Request.Context(), sectionID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "get questions by section success", result)
}

// FindByID godoc
//
// @Summary Get question by ID
// @Description Get question detail by ID
// @Tags Questions
// @Produce json
// @Param id path string true "Question ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /questions/{id} [get]
func (h *Handler) FindByID(c *gin.Context) {
	id := c.Param("id")

	result, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "get question success", result)
}

// Create godoc
//
// @Summary Create question
// @Description Create a new question
// @Tags Questions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body CreateQuestionRequest true "Create question request"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /questions [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateQuestionRequest

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

	response.Success(c, http.StatusCreated, "create question success", result)
}

// Update godoc
//
// @Summary Update question
// @Description Update existing question
// @Tags Questions
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Question ID"
// @Param request body UpdateQuestionRequest true "Update question request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /questions/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateQuestionRequest

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

	response.Success(c, http.StatusOK, "update question success", result)
}

// Delete godoc
//
// @Summary Delete question
// @Description Delete question by ID
// @Tags Questions
// @Security BearerAuth
// @Produce json
// @Param id path string true "Question ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /questions/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "delete question success", nil)
}
