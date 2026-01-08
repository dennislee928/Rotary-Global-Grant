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

// TrainingHandler handles training event HTTP requests
type TrainingHandler struct {
	trainingSvc *service.TrainingService
}

// NewTrainingHandler creates a new training handler
func NewTrainingHandler(trainingSvc *service.TrainingService) *TrainingHandler {
	return &TrainingHandler{trainingSvc: trainingSvc}
}

// Create handles POST /v1/training-events
// @Summary Create a training event
// @Description Create a new training event
// @Tags training
// @Accept json
// @Produce json
// @Param request body dto.CreateTrainingEventRequest true "Training event data"
// @Success 201 {object} vo.TrainingEventVO
// @Failure 400 {object} vo.ErrorVO
// @Failure 500 {object} vo.ErrorVO
// @Security BearerAuth
// @Router /v1/training-events [post]
func (h *TrainingHandler) Create(c *gin.Context) {
	var req dto.CreateTrainingEventRequest
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
	event, err := h.trainingSvc.CreateEvent(c.Request.Context(), req, userID, actorIP)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorVO{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// List handles GET /v1/training-events
// @Summary List training events
// @Description Get a paginated list of training events
// @Tags training
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20)
// @Param from query string false "From date (YYYY-MM-DD)"
// @Param to query string false "To date (YYYY-MM-DD)"
// @Success 200 {object} vo.TrainingEventListVO
// @Failure 400 {object} vo.ErrorVO
// @Failure 500 {object} vo.ErrorVO
// @Security BearerAuth
// @Router /v1/training-events [get]
func (h *TrainingHandler) List(c *gin.Context) {
	var query dto.ListTrainingEventsQuery
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

	events, err := h.trainingSvc.ListEvents(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorVO{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, events)
}

// GetByID handles GET /v1/training-events/:id
// @Summary Get training event by ID
// @Description Get detailed information about a specific training event
// @Tags training
// @Accept json
// @Produce json
// @Param id path string true "Training event ID"
// @Success 200 {object} vo.TrainingEventVO
// @Failure 404 {object} vo.ErrorVO
// @Failure 500 {object} vo.ErrorVO
// @Security BearerAuth
// @Router /v1/training-events/{id} [get]
func (h *TrainingHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	event, err := h.trainingSvc.GetEventByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrTrainingEventNotFound) {
			c.JSON(http.StatusNotFound, vo.ErrorVO{
				Code:    "NOT_FOUND",
				Message: "Training event not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, vo.ErrorVO{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get training event",
		})
		return
	}

	c.JSON(http.StatusOK, event)
}

// RecordQuizResult handles POST /v1/training-events/:id/results
// @Summary Record quiz result
// @Description Record a quiz result for a training event
// @Tags training
// @Accept json
// @Produce json
// @Param id path string true "Training event ID"
// @Param request body dto.RecordQuizResultRequest true "Quiz result data"
// @Success 201 {object} vo.QuizResultVO
// @Failure 400 {object} vo.ErrorVO
// @Failure 404 {object} vo.ErrorVO
// @Failure 500 {object} vo.ErrorVO
// @Security BearerAuth
// @Router /v1/training-events/{id}/results [post]
func (h *TrainingHandler) RecordQuizResult(c *gin.Context) {
	eventID := c.Param("id")

	var req dto.RecordQuizResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.ErrorVO{
			Code:    "VALIDATION_ERROR",
			Message: err.Error(),
		})
		return
	}

	result, err := h.trainingSvc.RecordQuizResult(c.Request.Context(), eventID, req)
	if err != nil {
		if errors.Is(err, service.ErrTrainingEventNotFound) {
			c.JSON(http.StatusNotFound, vo.ErrorVO{
				Code:    "NOT_FOUND",
				Message: "Training event not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, vo.ErrorVO{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to record quiz result",
		})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetStats handles GET /v1/training-events/stats
// @Summary Get training statistics
// @Description Get aggregated training statistics
// @Tags training
// @Accept json
// @Produce json
// @Success 200 {object} vo.TrainingStatsVO
// @Failure 500 {object} vo.ErrorVO
// @Security BearerAuth
// @Router /v1/training-events/stats [get]
func (h *TrainingHandler) GetStats(c *gin.Context) {
	stats, err := h.trainingSvc.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, vo.ErrorVO{
			Code:    "INTERNAL_ERROR",
			Message: "Failed to get training stats",
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}
