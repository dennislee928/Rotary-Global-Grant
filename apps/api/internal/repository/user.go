package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/dennislee928/rotary-global-grant-safety-resilience/apps/api/internal/model"
)

// UserRepository handles user database operations
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{db: db.Gorm}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete deletes a user
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, "id = ?", id).Error
}

// List retrieves all users with pagination
func (r *UserRepository) List(ctx context.Context, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	if err := r.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&users).Error

	return users, total, err
}

// CountByRole counts users by role
func (r *UserRepository) CountByRole(ctx context.Context, role string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("role = ? AND is_active = true", role).
		Count(&count).Error
	return count, err
}

// CreateAPIKey creates a new API key for a user
func (r *UserRepository) CreateAPIKey(ctx context.Context, apiKey *model.APIKey) error {
	return r.db.WithContext(ctx).Create(apiKey).Error
}

// GetAPIKeyByHash retrieves an API key by its hash
func (r *UserRepository) GetAPIKeyByHash(ctx context.Context, keyHash string) (*model.APIKey, error) {
	var apiKey model.APIKey
	err := r.db.WithContext(ctx).
		Preload("User").
		First(&apiKey, "key_hash = ?", keyHash).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &apiKey, err
}

// GetAPIKeysByUser retrieves all API keys for a user
func (r *UserRepository) GetAPIKeysByUser(ctx context.Context, userID uuid.UUID) ([]model.APIKey, error) {
	var keys []model.APIKey
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&keys).Error
	return keys, err
}

// UpdateAPIKeyLastUsed updates the last used timestamp of an API key
func (r *UserRepository) UpdateAPIKeyLastUsed(ctx context.Context, id uuid.UUID) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).
		Model(&model.APIKey{}).
		Where("id = ?", id).
		Update("last_used_at", &now).Error
}

// DeleteAPIKey deletes an API key
func (r *UserRepository) DeleteAPIKey(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.APIKey{}, "id = ?", id).Error
}

// HashAPIKey hashes an API key
func HashAPIKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}
