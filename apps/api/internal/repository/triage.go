package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/model"
)

// TriageRepository handles triage decision database operations
type TriageRepository struct {
	db *gorm.DB
}

// NewTriageRepository creates a new triage repository
func NewTriageRepository(db *DB) *TriageRepository {
	return &TriageRepository{db: db.Gorm}
}

// Create creates a new triage decision
func (r *TriageRepository) Create(ctx context.Context, decision *model.TriageDecision) error {
	return r.db.WithContext(ctx).Create(decision).Error
}

// GetByID retrieves a triage decision by ID
func (r *TriageRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.TriageDecision, error) {
	var decision model.TriageDecision
	err := r.db.WithContext(ctx).
		Preload("Decider").
		Preload("Report").
		First(&decision, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &decision, err
}

// GetByReportID retrieves all triage decisions for a report
func (r *TriageRepository) GetByReportID(ctx context.Context, reportID uuid.UUID) ([]model.TriageDecision, error) {
	var decisions []model.TriageDecision
	err := r.db.WithContext(ctx).
		Preload("Decider").
		Where("report_id = ?", reportID).
		Order("decided_at DESC").
		Find(&decisions).Error
	return decisions, err
}

// List retrieves triage decisions with pagination
func (r *TriageRepository) List(ctx context.Context, params ListTriageParams) ([]model.TriageDecision, int64, error) {
	var decisions []model.TriageDecision
	var total int64

	query := r.db.WithContext(ctx).Model(&model.TriageDecision{})

	// Apply filters
	if params.ReportID != uuid.Nil {
		query = query.Where("report_id = ?", params.ReportID)
	}
	if params.Decision != "" {
		query = query.Where("decision = ?", params.Decision)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	query = query.
		Preload("Decider").
		Preload("Report").
		Order("decided_at DESC").
		Offset(offset).
		Limit(params.PageSize)

	err := query.Find(&decisions).Error
	return decisions, total, err
}

// GetLatestForReport retrieves the latest triage decision for a report
func (r *TriageRepository) GetLatestForReport(ctx context.Context, reportID uuid.UUID) (*model.TriageDecision, error) {
	var decision model.TriageDecision
	err := r.db.WithContext(ctx).
		Preload("Decider").
		Where("report_id = ?", reportID).
		Order("decided_at DESC").
		First(&decision).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &decision, err
}

// GetMedianTriageTime calculates median time from report creation to first triage
func (r *TriageRepository) GetMedianTriageTime(ctx context.Context) (time.Duration, error) {
	var result struct {
		MedianMinutes float64
	}

	err := r.db.WithContext(ctx).Raw(`
		WITH first_triage AS (
			SELECT 
				report_id,
				MIN(decided_at) as first_decision
			FROM triage_decisions
			GROUP BY report_id
		)
		SELECT 
			PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY 
				EXTRACT(EPOCH FROM (ft.first_decision - r.created_at)) / 60
			) as median_minutes
		FROM first_triage ft
		JOIN reports r ON r.id = ft.report_id
	`).Scan(&result).Error

	if err != nil {
		return 0, err
	}

	return time.Duration(result.MedianMinutes) * time.Minute, nil
}

// GetVerifiedRatio calculates the ratio of accepted (verified) reports
func (r *TriageRepository) GetVerifiedRatio(ctx context.Context) (float64, error) {
	var result struct {
		Ratio float64
	}

	err := r.db.WithContext(ctx).Raw(`
		WITH latest_decisions AS (
			SELECT DISTINCT ON (report_id)
				report_id, decision
			FROM triage_decisions
			ORDER BY report_id, decided_at DESC
		)
		SELECT 
			COALESCE(
				SUM(CASE WHEN decision = 'accept' THEN 1 ELSE 0 END)::float / 
				NULLIF(COUNT(*), 0), 0
			) * 100 as ratio
		FROM latest_decisions
	`).Scan(&result).Error

	return result.Ratio, err
}

// GetAbuseRate calculates the abuse/spam report rate
func (r *TriageRepository) GetAbuseRate(ctx context.Context) (float64, error) {
	var result struct {
		Rate float64
	}

	err := r.db.WithContext(ctx).Raw(`
		SELECT 
			COALESCE(
				SUM(CASE WHEN status = 'spam' THEN 1 ELSE 0 END)::float / 
				NULLIF(COUNT(*), 0), 0
			) * 100 as rate
		FROM reports
	`).Scan(&result).Error

	return result.Rate, err
}

// ListTriageParams represents parameters for listing triage decisions
type ListTriageParams struct {
	Page     int
	PageSize int
	ReportID uuid.UUID
	Decision string
}
