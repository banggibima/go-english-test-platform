package auth

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

// Register godoc
//
// @Summary Register user
// @Description Create a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register request"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if err := validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, "validation failed", err.Error())
		return
	}

	result, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusCreated, "register success", result)
}

// Login godoc
//
// @Summary Login user
// @Description Authenticate user and return access token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if err := validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, "validation failed", err.Error())
		return
	}

	result, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "login success", result)
}

// RefreshToken godoc
//
// @Summary Refresh access token
// @Description Generate a new access token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body", err.Error())
		return
	}

	if err := validator.Validate(req); err != nil {
		response.Error(c, http.StatusBadRequest, "validation failed", err.Error())
		return
	}

	result, err := h.service.RefreshToken(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	response.Success(c, http.StatusOK, "refresh token success", result)
}
