package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/model"
)

// AuditRepository handles audit log database operations
type AuditRepository struct {
	db *gorm.DB
}

// NewAuditRepository creates a new audit repository
func NewAuditRepository(db *DB) *AuditRepository {
	return &AuditRepository{db: db.Gorm}
}

// Create creates a new audit log entry
func (r *AuditRepository) Create(ctx context.Context, log *model.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// List retrieves audit logs with pagination
func (r *AuditRepository) List(ctx context.Context, params ListAuditParams) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	query := r.db.WithContext(ctx).Model(&model.AuditLog{})

	// Apply filters
	if params.ObjectType != "" {
		query = query.Where("object_type = ?", params.ObjectType)
	}
	if params.ObjectID != uuid.Nil {
		query = query.Where("object_id = ?", params.ObjectID)
	}
	if params.ActorID != uuid.Nil {
		query = query.Where("actor_id = ?", params.ActorID)
	}
	if params.Action != "" {
		query = query.Where("action = ?", params.Action)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	query = query.
		Preload("Actor").
		Order("ts DESC").
		Offset(offset).
		Limit(params.PageSize)

	err := query.Find(&logs).Error
	return logs, total, err
}

// GetByObjectID retrieves audit logs for a specific object
func (r *AuditRepository) GetByObjectID(ctx context.Context, objectType string, objectID uuid.UUID) ([]model.AuditLog, error) {
	var logs []model.AuditLog
	err := r.db.WithContext(ctx).
		Preload("Actor").
		Where("object_type = ? AND object_id = ?", objectType, objectID).
		Order("ts DESC").
		Find(&logs).Error
	return logs, err
}

// ListAuditParams represents parameters for listing audit logs
type ListAuditParams struct {
	Page       int
	PageSize   int
	ObjectType string
	ObjectID   uuid.UUID
	ActorID    uuid.UUID
	Action     string
}
