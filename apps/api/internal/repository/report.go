package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/model"
)

// ReportRepository handles report database operations
type ReportRepository struct {
	db *gorm.DB
}

// NewReportRepository creates a new report repository
func NewReportRepository(db *DB) *ReportRepository {
	return &ReportRepository{db: db.Gorm}
}

// Create creates a new report
func (r *ReportRepository) Create(ctx context.Context, report *model.Report) error {
	return r.db.WithContext(ctx).Create(report).Error
}

// GetByID retrieves a report by ID
func (r *ReportRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Report, error) {
	var report model.Report
	err := r.db.WithContext(ctx).
		Preload("TriageDecisions").
		First(&report, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &report, err
}

// List retrieves reports with pagination and filtering
func (r *ReportRepository) List(ctx context.Context, params ListReportParams) ([]model.Report, int64, error) {
	var reports []model.Report
	var total int64

	query := r.db.WithContext(ctx).Model(&model.Report{})

	// Apply filters
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.Category != "" {
		query = query.Where("category = ?", params.Category)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	orderBy := "created_at DESC"
	if params.SortBy != "" {
		dir := "DESC"
		if params.SortDir == "asc" {
			dir = "ASC"
		}
		orderBy = params.SortBy + " " + dir
	}
	query = query.Order(orderBy)

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	err := query.Find(&reports).Error
	return reports, total, err
}

// Update updates a report
func (r *ReportRepository) Update(ctx context.Context, report *model.Report) error {
	return r.db.WithContext(ctx).Save(report).Error
}

// UpdateStatus updates the status of a report
func (r *ReportRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	return r.db.WithContext(ctx).
		Model(&model.Report{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// Delete deletes a report
func (r *ReportRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Report{}, "id = ?", id).Error
}

// GetStats retrieves report statistics
func (r *ReportRepository) GetStats(ctx context.Context) (*ReportStats, error) {
	var stats ReportStats

	// Total reports
	if err := r.db.WithContext(ctx).Model(&model.Report{}).Count(&stats.Total).Error; err != nil {
		return nil, err
	}

	// Reports by status
	var statusCounts []struct {
		Status string
		Count  int64
	}
	if err := r.db.WithContext(ctx).
		Model(&model.Report{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&statusCounts).Error; err != nil {
		return nil, err
	}
	stats.ByStatus = make(map[string]int64)
	for _, sc := range statusCounts {
		stats.ByStatus[sc.Status] = sc.Count
	}

	// Reports by category
	var categoryCounts []struct {
		Category string
		Count    int64
	}
	if err := r.db.WithContext(ctx).
		Model(&model.Report{}).
		Select("category, count(*) as count").
		Group("category").
		Scan(&categoryCounts).Error; err != nil {
		return nil, err
	}
	stats.ByCategory = make(map[string]int64)
	for _, cc := range categoryCounts {
		stats.ByCategory[cc.Category] = cc.Count
	}

	return &stats, nil
}

// ListReportParams represents parameters for listing reports
type ListReportParams struct {
	Page     int
	PageSize int
	Status   string
	Category string
	SortBy   string
	SortDir  string
}

// ReportStats represents report statistics
type ReportStats struct {
	Total      int64
	ByStatus   map[string]int64
	ByCategory map[string]int64
}
