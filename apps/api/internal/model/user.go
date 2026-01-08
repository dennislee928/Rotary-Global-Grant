package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a system user (triager, admin, auditor, educator)
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Email        string    `gorm:"uniqueIndex;size:255;not null"`
	PasswordHash string    `gorm:"size:255;not null"`
	Role         string    `gorm:"size:50;not null;default:'triager'"`
	DisplayName  string    `gorm:"size:255"`
	IsActive     bool      `gorm:"not null;default:true"`
	CreatedAt    time.Time `gorm:"not null;default:now()"`
	UpdatedAt    time.Time `gorm:"not null;default:now()"`

	// Associations
	TriageDecisions []TriageDecision `gorm:"foreignKey:DecidedBy"`
	Alerts          []Alert          `gorm:"foreignKey:ApprovedBy"`
	AuditLogs       []AuditLog       `gorm:"foreignKey:ActorID"`
	APIKeys         []APIKey         `gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// UserRole constants
const (
	RoleAdmin    = "admin"
	RoleTriager  = "triager"
	RoleAuditor  = "auditor"
	RoleEducator = "educator"
)

// ValidRoles returns all valid user roles
func ValidRoles() []string {
	return []string{RoleAdmin, RoleTriager, RoleAuditor, RoleEducator}
}
