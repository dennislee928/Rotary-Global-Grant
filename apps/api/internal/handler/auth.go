package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/dto"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/service"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/vo"
)

// AuthHandler handles authentication HTTP requests
type AuthHandler struct {
	authSvc *service.AuthService
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// Login handles POST /v1/auth/login
// @Summary User login
// @Description Authenticate user and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} vo.TokenVO
// @Failure 400 {object} vo.ErrorVO
// @Failure 401 {object} vo.ErrorVO
// @Router /v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorVO{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	actorIP := c.ClientIP()
	token, err := h.authSvc.Login(c.Request.Context(), req, actorIP)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, vo.ErrorVO{
				Code:    "UNAUTHORIZED",
				Message: "Invalid email or password",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, vo.ErrorVO{
			Code:    "INTERNAL_ERROR",
			Message: "Login failed",
		})
		return
	}

	c.JSON(http.StatusOK, token)
}

// Me handles GET /v1/auth/me
// @Summary Get current user
// @Description Get information about the currently authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} vo.UserSummaryVO
// @Failure 401 {object} vo.ErrorVO
// @Security BearerAuth
// @Router /v1/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, vo.ErrorVO{
			Code:    "UNAUTHORIZED",
			Message: "Not authenticated",
		})
		return
	}

	user, err := h.authSvc.GetUserByID(c.Request.Context(), userID.(uuid.UUID).String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorVO{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get user",
		})
		return
	}

	c.JSON(http.StatusOK, vo.UserSummaryVO{
		ID:          user.ID.String(),
		DisplayName: user.DisplayName,
		Role:        user.Role,
	})
}

// CreateAPIKey handles POST /v1/auth/api-keys
// @Summary Create API key
// @Description Create a new API key for the current user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.CreateAPIKeyRequest true "API key data"
// @Success 201 {object} vo.APIKeyVO
// @Failure 400 {object} vo.ErrorVO
// @Failure 401 {object} vo.ErrorVO
// @Failure 500 {object} vo.ErrorVO
// @Security BearerAuth
// @Router /v1/auth/api-keys [post]
func (h *AuthHandler) CreateAPIKey(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, vo.ErrorVO{
			Code:    "UNAUTHORIZED",
			Message: "Not authenticated",
		})
		return
	}

	var req dto.CreateAPIKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorVO{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	actorIP := c.ClientIP()
	apiKey, err := h.authSvc.CreateAPIKey(c.Request.Context(), userID.(uuid.UUID), req, actorIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorVO{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to create API key",
		})
		return
	}

	c.JSON(http.StatusCreated, apiKey)
}
