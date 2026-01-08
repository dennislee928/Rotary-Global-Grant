package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/dto"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/service"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/vo"
)

// ReportHandler handles report HTTP requests
type ReportHandler struct {
	reportSvc *service.ReportService
}

// NewReportHandler creates a new report handler
func NewReportHandler(reportSvc *service.ReportService) *ReportHandler {
	return &ReportHandler{reportSvc: reportSvc}
}

// Create handles POST /v1/reports
// @Summary Create a new report
// @Description Submit a new community incident report
// @Tags reports
// @Accept json
// @Produce json
// @Param request body dto.CreateReportRequest true "Report data"
// @Success 201 {object} vo.ReportVO
// @Failure 400 {object} vo.ErrorVO
// @Failure 500 {object} vo.ErrorVO
// @Router /v1/reports [post]
func (h *ReportHandler) Create(c *gin.Context) {
	var req dto.CreateReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorVO{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	actorIP := c.ClientIP()
	report, err := h.reportSvc.Create(c.Request.Context(), req, actorIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorVO{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to create report",
		})
		return
	}

	c.JSON(http.StatusCreated, report)
}

// List handles GET /v1/reports
// @Summary List reports
// @Description Get a paginated list of reports
// @Tags reports
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Param status query string false "Filter by status"
// @Param category query string false "Filter by category"
// @Param sortBy query string false "Sort by field" default(createdAt)
// @Param sortDir query string false "Sort direction" default(desc)
// @Success 200 {object} vo.ReportListVO
// @Failure 400 {object} vo.ErrorVO
// @Failure 500 {object} vo.ErrorVO
// @Security BearerAuth
// @Router /v1/reports [get]
func (h *ReportHandler) List(c *gin.Context) {
	var query dto.ListReportsQuery
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
	if query.SortBy == "" {
		query.SortBy = "createdAt"
	}
	if query.SortDir == "" {
		query.SortDir = "desc"
	}

	reports, err := h.reportSvc.List(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorVO{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to list reports",
		})
		return
	}

	c.JSON(http.StatusOK, reports)
}

// GetByID handles GET /v1/reports/:id
// @Summary Get report by ID
// @Description Get detailed information about a specific report
// @Tags reports
// @Accept json
// @Produce json
// @Param id path string true "Report ID"
// @Success 200 {object} vo.ReportDetailVO
// @Failure 404 {object} vo.ErrorVO
// @Failure 500 {object} vo.ErrorVO
// @Security BearerAuth
// @Router /v1/reports/{id} [get]
func (h *ReportHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	report, err := h.reportSvc.GetByID(c.Request.Context(), id)
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
			Message: "Failed to get report",
		})
		return
	}

	c.JSON(http.StatusOK, report)
}
