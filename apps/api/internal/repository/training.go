package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/model"
)

// TrainingRepository handles training event database operations
type TrainingRepository struct {
	db *gorm.DB
}

// NewTrainingRepository creates a new training repository
func NewTrainingRepository(db *DB) *TrainingRepository {
	return &TrainingRepository{db: db.Gorm}
}

// CreateEvent creates a new training event
func (r *TrainingRepository) CreateEvent(ctx context.Context, event *model.TrainingEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

// GetEventByID retrieves a training event by ID
func (r *TrainingRepository) GetEventByID(ctx context.Context, id uuid.UUID) (*model.TrainingEvent, error) {
	var event model.TrainingEvent
	err := r.db.WithContext(ctx).
		Preload("Creator").
		First(&event, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &event, err
}

// ListEvents retrieves training events with pagination
func (r *TrainingRepository) ListEvents(ctx context.Context, params ListTrainingParams) ([]model.TrainingEvent, int64, error) {
	var events []model.TrainingEvent
	var total int64

	query := r.db.WithContext(ctx).Model(&model.TrainingEvent{})

	// Apply date filters
	if !params.From.IsZero() {
		query = query.Where("event_date >= ?", params.From)
	}
	if !params.To.IsZero() {
		query = query.Where("event_date <= ?", params.To)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	query = query.
		Preload("Creator").
		Order("event_date DESC").
		Offset(offset).
		Limit(params.PageSize)

	err := query.Find(&events).Error
	return events, total, err
}

// UpdateEvent updates a training event
func (r *TrainingRepository) UpdateEvent(ctx context.Context, event *model.TrainingEvent) error {
	return r.db.WithContext(ctx).Save(event).Error
}

// DeleteEvent deletes a training event
func (r *TrainingRepository) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.TrainingEvent{}, "id = ?", id).Error
}

// CreateQuizResult creates a quiz result
func (r *TrainingRepository) CreateQuizResult(ctx context.Context, result *model.QuizResult) error {
	return r.db.WithContext(ctx).Create(result).Error
}

// GetQuizResultsByEvent retrieves quiz results for an event
func (r *TrainingRepository) GetQuizResultsByEvent(ctx context.Context, eventID uuid.UUID) ([]model.QuizResult, error) {
	var results []model.QuizResult
	err := r.db.WithContext(ctx).
		Where("event_id = ?", eventID).
		Order("created_at DESC").
		Find(&results).Error
	return results, err
}

// AddParticipant adds a participant to an event
func (r *TrainingRepository) AddParticipant(ctx context.Context, participant *model.TrainingParticipant) error {
	return r.db.WithContext(ctx).Create(participant).Error
}

// GetTrainingStats retrieves training statistics
func (r *TrainingRepository) GetTrainingStats(ctx context.Context) (*TrainingStats, error) {
	var stats TrainingStats

	// Total events
	if err := r.db.WithContext(ctx).Model(&model.TrainingEvent{}).Count(&stats.TotalEvents).Error; err != nil {
		return nil, err
	}

	// Total participants (sum of attendance counts)
	var result struct {
		Total int64
	}
	if err := r.db.WithContext(ctx).
		Model(&model.TrainingEvent{}).
		Select("COALESCE(SUM(attendance_count), 0) as total").
		Scan(&result).Error; err != nil {
		return nil, err
	}
	stats.TotalParticipants = result.Total

	// Average improvement
	var avgResult struct {
		AvgImprovement float64
	}
	if err := r.db.WithContext(ctx).Raw(`
		SELECT COALESCE(AVG(post_avg - pre_avg), 0) as avg_improvement
		FROM training_events
		WHERE pre_avg IS NOT NULL AND post_avg IS NOT NULL
	`).Scan(&avgResult).Error; err != nil {
		return nil, err
	}
	stats.AverageImprovement = avgResult.AvgImprovement

	return &stats, nil
}

// ListTrainingParams represents parameters for listing training events
type ListTrainingParams struct {
	Page     int
	PageSize int
	From     time.Time
	To       time.Time
}

// TrainingStats represents training statistics
type TrainingStats struct {
	TotalEvents        int64
	TotalParticipants  int64
	AverageImprovement float64
}
