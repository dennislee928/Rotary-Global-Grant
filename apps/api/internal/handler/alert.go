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

// AlertHandler handles alert HTTP requests
type AlertHandler struct {
	alertSvc *service.AlertService
}

// NewAlertHandler creates a new alert handler
func NewAlertHandler(alertSvc *service.AlertService) *AlertHandler {
	return &AlertHandler{alertSvc: alertSvc}
}

// Create handles POST /v1/alerts
// @Summary Create an alert
// @Description Create a new CAP-ready alert
// @Tags alerts
// @Accept json
// @Produce json
// @Param request body dto.CreateAlertRequest true "Alert data"
// @Success 201 {object} vo.AlertVO
// @Failure 400 {object} vo.ErrorVO
// @Failure 500 {object} vo.ErrorVO
// @Security BearerAuth
// @Router /v1/alerts [post]
func (h *AlertHandler) Create(c *gin.Context) {
	var req dto.CreateAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorVO{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	var userID *uuid.UUID
	if uid, exists := c.Get("userID"); exists {
		if id, ok := uid.(uuid.UUID); ok {
			userID = &id
		}
	}

	actorIP := c.ClientIP()
	alert, err := h.alertSvc.Create(c.Request.Context(), req, userID, actorIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorVO{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to create alert",
		})
		return
	}

	c.JSON(http.StatusCreated, alert)
}

// List handles GET /v1/alerts
// @Summary List alerts
// @Description Get a paginated list of alerts
// @Tags alerts
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Param status query string false "Filter by status"
// @Success 200 {object} vo.AlertListVO
// @Failure 400 {object} vo.ErrorVO
// @Failure 500 {object} vo.ErrorVO
// @Security BearerAuth
// @Router /v1/alerts [get]
func (h *AlertHandler) List(c *gin.Context) {
	var query dto.ListAlertsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorVO{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}

	alerts, err := h.alertSvc.List(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorVO{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to list alerts",
		})
		return
	}

	c.JSON(http.StatusOK, alerts)
}

// GetByID handles GET /v1/alerts/:id
// @Summary Get alert by ID
// @Description Get detailed information about a specific alert
// @Tags alerts
// @Accept json
// @Produce json
// @Param id path string true "Alert ID"
// @Success 200 {object} vo.AlertVO
// @Failure 404 {object} vo.ErrorVO
// @Failure 500 {object} vo.ErrorVO
// @Security BearerAuth
// @Router /v1/alerts/{id} [get]
func (h *AlertHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	alert, err := h.alertSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrAlertNotFound) {
			c.JSON(http.StatusNotFound, vo.ErrorVO{
				Code:    "NOT_FOUND",
				Message: "Alert not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, vo.ErrorVO{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get alert",
		})
		return
	}

	c.JSON(http.StatusOK, alert)
}

// Update handles PATCH /v1/alerts/:id
// @Summary Update an alert
// @Description Update alert status or content
// @Tags alerts
// @Accept json
// @Produce json
// @Param id path string true "Alert ID"
// @Param request body dto.UpdateAlertRequest true "Update data"
// @Success 200 {object} vo.AlertVO
// @Failure 400 {object} vo.ErrorVO
// @Failure 404 {object} vo.ErrorVO
// @Failure 500 {object} vo.ErrorVO
// @Security BearerAuth
// @Router /v1/alerts/{id} [patch]
func (h *AlertHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateAlertRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorVO{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	var userID *uuid.UUID
	if uid, exists := c.Get("userID"); exists {
		if id, ok := uid.(uuid.UUID); ok {
			userID = &id
		}
	}

	actorIP := c.ClientIP()
	alert, err := h.alertSvc.Update(c.Request.Context(), id, req, userID, actorIP)
	if err != nil {
		if errors.Is(err, service.ErrAlertNotFound) {
			c.JSON(http.StatusNotFound, vo.ErrorVO{
				Code:    "NOT_FOUND",
				Message: "Alert not found",
			})
			return
		}
		if errors.Is(err, service.ErrInvalidTransition) {
			c.JSON(http.StatusBadRequest, vo.ErrorVO{
				Code:    "INVALID_TRANSITION",
				Message: "Invalid status transition",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, vo.ErrorVO{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to update alert",
		})
		return
	}

	c.JSON(http.StatusOK, alert)
}
