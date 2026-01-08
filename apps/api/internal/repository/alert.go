package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/model"
)

// AlertRepository handles alert database operations
type AlertRepository struct {
	db *gorm.DB
}

// NewAlertRepository creates a new alert repository
func NewAlertRepository(db *DB) *AlertRepository {
	return &AlertRepository{db: db.Gorm}
}

// Create creates a new alert
func (r *AlertRepository) Create(ctx context.Context, alert *model.Alert) error {
	return r.db.WithContext(ctx).Create(alert).Error
}

// GetByID retrieves an alert by ID
func (r *AlertRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Alert, error) {
	var alert model.Alert
	err := r.db.WithContext(ctx).
		Preload("Approver").
		Preload("Report").
		First(&alert, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &alert, err
}

// List retrieves alerts with pagination and filtering
func (r *AlertRepository) List(ctx context.Context, params ListAlertParams) ([]model.Alert, int64, error) {
	var alerts []model.Alert
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Alert{})

	// Apply filters
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	query = query.
		Preload("Approver").
		Order("created_at DESC").
		Offset(offset).
		Limit(params.PageSize)

	err := query.Find(&alerts).Error
	return alerts, total, err
}

// Update updates an alert
func (r *AlertRepository) Update(ctx context.Context, alert *model.Alert) error {
	return r.db.WithContext(ctx).Save(alert).Error
}

// UpdateStatus updates alert status with optional approver and publish time
func (r *AlertRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string, approverID *uuid.UUID) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now().UTC(),
	}

	if approverID != nil {
		updates["approved_by"] = approverID
	}

	if status == model.AlertStatusPublished {
		now := time.Now().UTC()
		updates["published_at"] = &now
	}

	return r.db.WithContext(ctx).
		Model(&model.Alert{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// GetActiveAlerts retrieves currently active (published) alerts
func (r *AlertRepository) GetActiveAlerts(ctx context.Context, limit int) ([]model.Alert, error) {
	var alerts []model.Alert
	err := r.db.WithContext(ctx).
		Where("status = ?", model.AlertStatusPublished).
		Order("published_at DESC").
		Limit(limit).
		Find(&alerts).Error
	return alerts, err
}

// GetRecentAlerts retrieves recent alerts regardless of status
func (r *AlertRepository) GetRecentAlerts(ctx context.Context, limit int) ([]model.Alert, error) {
	var alerts []model.Alert
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(limit).
		Find(&alerts).Error
	return alerts, err
}

// GetPublishLatency calculates median time from approval to publish
func (r *AlertRepository) GetPublishLatency(ctx context.Context) (time.Duration, error) {
	var result struct {
		MedianMinutes float64
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			COALESCE(
				PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY 
					EXTRACT(EPOCH FROM (published_at - created_at)) / 60
				), 0
			) as median_minutes
		FROM alerts
		WHERE status = 'published' AND published_at IS NOT NULL
	`).Scan(&result).Error

	if err != nil {
		return 0, err
	}

	return time.Duration(result.MedianMinutes) * time.Minute, nil
}

// CountByStatus counts alerts by status
func (r *AlertRepository) CountByStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Alert{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}

// ListAlertParams represents parameters for listing alerts
type ListAlertParams struct {
	Page     int
	PageSize int
	Status   string
}
