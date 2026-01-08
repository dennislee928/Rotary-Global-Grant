package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/service"
	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/vo"
)

// MetricsHandler handles metrics HTTP requests
type MetricsHandler struct {
	metricsSvc *service.MetricsService
}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler(metricsSvc *service.MetricsService) *MetricsHandler {
	return &MetricsHandler{metricsSvc: metricsSvc}
}

// GetKPIMetrics handles GET /v1/metrics/kpi
// @Summary Get KPI metrics
// @Description Get all KPI metrics for the project
// @Tags metrics
// @Accept json
// @Produce json
// @Success 200 {object} vo.KPIMetricsVO
// @Failure 500 {object} vo.ErrorVO
// @Security BearerAuth
// @Router /v1/metrics/kpi [get]
func (h *MetricsHandler) GetKPIMetrics(c *gin.Context) {
	metrics, err := h.metricsSvc.GetKPIMetrics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorVO{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get KPI metrics",
		})
		return
	}

	c.JSON(http.StatusOK, metrics)
}

// GetDashboardStats handles GET /v1/metrics/dashboard
// @Summary Get dashboard statistics
// @Description Get public dashboard statistics
// @Tags metrics
// @Accept json
// @Produce json
// @Success 200 {object} vo.DashboardStatsVO
// @Failure 500 {object} vo.ErrorVO
// @Router /v1/metrics/dashboard [get]
func (h *MetricsHandler) GetDashboardStats(c *gin.Context) {
	stats, err := h.metricsSvc.GetDashboardStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorVO{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get dashboard stats",
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}
