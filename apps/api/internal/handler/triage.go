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

// TriageHandler handles triage HTTP requests
type TriageHandler struct {
	triageSvc *service.TriageService
}

// NewTriageHandler creates a new triage handler
func NewTriageHandler(triageSvc *service.TriageService) *TriageHandler {
	return &TriageHandler{triageSvc: triageSvc}
}

// TriageReport handles POST /v1/reports/:id/triage
// @Summary Triage a report
// @Description Create a triage decision for a report
// @Tags triage
// @Accept json
// @Produce json
// @Param id path string true "Report ID"
// @Param request body dto.TriageRequest true "Triage decision"
// @Success 200 {object} vo.TriageDecisionVO
// @Failure 400 {object} vo.ErrorVO
// @Failure 404 {object} vo.ErrorVO
// @Failure 500 {object} vo.ErrorVO
// @Security BearerAuth
// @Router /v1/reports/{id}/triage [post]
func (h *TriageHandler) TriageReport(c *gin.Context) {
	reportID := c.Param("id")

	var req dto.TriageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorVO{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	var userID *uuid.UUID
	if uid, exists := c.Get("userID"); exists {
		if id, ok := uid.(uuid.UUID); ok {
			userID = &id
		}
	}

	actorIP := c.ClientIP()
	decision, err := h.triageSvc.TriageReport(c.Request.Context(), reportID, req, userID, actorIP)
	if err != nil {
		if errors.Is(err, service.ErrReportNotFound) {
			c.JSON(http.StatusNotFound, vo.ErrorVO{
				Code:    "NOT_FOUND",
				Message: "Report not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, vo.ErrorVO{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to triage report",
		})
		return
	}

	c.JSON(http.StatusOK, decision)
}

// List handles GET /v1/triage-decisions
// @Summary List triage decisions
// @Description Get a paginated list of triage decisions
// @Tags triage
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Param reportId query string false "Filter by report ID"
// @Param decision query string false "Filter by decision type"
// @Success 200 {object} vo.TriageDecisionListVO
// @Failure 400 {object} vo.ErrorVO
// @Failure 500 {object} vo.ErrorVO
// @Security BearerAuth
// @Router /v1/triage-decisions [get]
func (h *TriageHandler) List(c *gin.Context) {
	var query dto.ListTriageDecisionsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorVO{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	// Set defaults
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}

	decisions, err := h.triageSvc.List(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorVO{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to list triage decisions",
		})
		return
	}

	c.JSON(http.StatusOK, decisions)
}
