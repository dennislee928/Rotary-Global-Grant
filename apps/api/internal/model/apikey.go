package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// APIKey represents an API key for authentication
type APIKey struct {
	ID         uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID     uuid.UUID   `gorm:"type:uuid;not null;index"`
	KeyHash    string      `gorm:"size:64;not null;uniqueIndex"`
	Name       string      `gorm:"size:255;not null"`
	Scopes     StringArray `gorm:"type:jsonb;default:'[\"read\"]'"`
	ExpiresAt  *time.Time
	LastUsedAt *time.Time
	CreatedAt  time.Time `gorm:"not null;default:now()"`
	IsActive   bool      `gorm:"not null;default:true"`

	// Associations
	User User `gorm:"foreignKey:UserID"`
}

func (APIKey) TableName() string {
	return "api_keys"
}

func (k *APIKey) BeforeCreate(tx *gorm.DB) error {
	if k.ID == uuid.Nil {
		k.ID = uuid.New()
	}
	return nil
}

// API key scopes
const (
	ScopeRead       = "read"
	ScopeWrite      = "write"
	ScopeTriage     = "triage"
	ScopeAlerts     = "alerts"
	ScopeAdmin      = "admin"
)

// ValidScopes returns all valid API key scopes
func ValidScopes() []string {
	return []string{ScopeRead, ScopeWrite, ScopeTriage, ScopeAlerts, ScopeAdmin}
}

// HasScope checks if the API key has a specific scope
func (k *APIKey) HasScope(scope string) bool {
	for _, s := range k.Scopes {
		if s == scope || s == ScopeAdmin {
			return true
		}
	}
	return false
}

// IsValid checks if the API key is valid (active and not expired)
func (k *APIKey) IsValid() bool {
	if !k.IsActive {
		return false
	}
	if k.ExpiresAt != nil && k.ExpiresAt.Before(time.Now()) {
		return false
	}
	return true
}
