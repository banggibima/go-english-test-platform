package files

import (
	"net/http"

	"github.com/banggibima/go-english-test-platform/pkg/response"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Upload godoc
//
// @Summary Upload file
// @Description Upload file to MinIO and store metadata
// @Tags Files
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /files/upload [post]
func (h *Handler) Upload(c *gin.Context) {
	userID := c.GetString("user_id")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "file is required", err.Error())
		return
	}

	result, err := h.service.Upload(c.Request.Context(), userID, fileHeader)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "upload file success", result)
}

// FindMyFiles godoc
//
// @Summary Get my files
// @Description Get uploaded files for current user
// @Tags Files
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /files [get]
func (h *Handler) FindMyFiles(c *gin.Context) {
	userID := c.GetString("user_id")

	result, err := h.service.FindByUserID(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "get files success", result)
}

// FindByID godoc
//
// @Summary Get file by ID
// @Description Get file metadata by ID
// @Tags Files
// @Security BearerAuth
// @Produce json
// @Param id path string true "File ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /files/{id} [get]
func (h *Handler) FindByID(c *gin.Context) {
	id := c.Param("id")

	result, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if result == nil {
		response.Error(c, http.StatusNotFound, "file not found", nil)
		return
	}

	response.Success(c, http.StatusOK, "get file success", result)
}
