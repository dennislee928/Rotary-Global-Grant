package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditLog represents an immutable audit log entry
type AuditLog struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ActorID    *uuid.UUID `gorm:"type:uuid;index"`
	ActorIP    string     `gorm:"size:45"`
	Action     string     `gorm:"size:100;not null"`
	ObjectType string     `gorm:"size:50;not null;index:idx_audit_object"`
	ObjectID   *uuid.UUID `gorm:"type:uuid;index:idx_audit_object"`
	Diff       JSONMap    `gorm:"type:jsonb"`
	Timestamp  time.Time  `gorm:"column:ts;not null;default:now();index"`

	// Associations
	Actor *User `gorm:"foreignKey:ActorID"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now().UTC()
	}
	return nil
}

// Audit actions
const (
	ActionCreate   = "create"
	ActionUpdate   = "update"
	ActionDelete   = "delete"
	ActionTriage   = "triage"
	ActionApprove  = "approve"
	ActionPublish  = "publish"
	ActionWithdraw = "withdraw"
	ActionLogin    = "login"
	ActionLogout   = "logout"
)

// Audit object types
const (
	ObjectTypeReport    = "report"
	ObjectTypeTriage    = "triage_decision"
	ObjectTypeAlert     = "alert"
	ObjectTypeTraining  = "training_event"
	ObjectTypeUser      = "user"
	ObjectTypeAPIKey    = "api_key"
)

// ValidAuditActions returns all valid audit actions
func ValidAuditActions() []string {
	return []string{
		ActionCreate, ActionUpdate, ActionDelete, ActionTriage,
		ActionApprove, ActionPublish, ActionWithdraw, ActionLogin, ActionLogout,
	}
}

// ValidObjectTypes returns all valid audit object types
func ValidObjectTypes() []string {
	return []string{
		ObjectTypeReport, ObjectTypeTriage, ObjectTypeAlert,
		ObjectTypeTraining, ObjectTypeUser, ObjectTypeAPIKey,
	}
}
