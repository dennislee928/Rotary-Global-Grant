package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TriageDecision represents a triage decision for a report
type TriageDecision struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ReportID      uuid.UUID  `gorm:"type:uuid;not null;index"`
	DecidedBy     *uuid.UUID `gorm:"type:uuid;index"`
	Decision      string     `gorm:"size:50;not null"`
	SeverityFinal string     `gorm:"size:10;not null"`
	EvidenceLevel string     `gorm:"size:10"`
	Rationale     string     `gorm:"type:text"`
	AuditHash     string     `gorm:"size:64"`
	DecidedAt     time.Time  `gorm:"not null;default:now()"`

	// Associations
	Report  Report `gorm:"foreignKey:ReportID"`
	Decider *User  `gorm:"foreignKey:DecidedBy"`
}

func (TriageDecision) TableName() string {
	return "triage_decisions"
}

func (t *TriageDecision) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// Triage decision types
const (
	DecisionAccept       = "accept"
	DecisionReject       = "reject"
	DecisionNeedsMoreInfo = "needs_more_info"
	DecisionEscalate     = "escalate"
)

// ValidDecisions returns all valid triage decisions
func ValidDecisions() []string {
	return []string{DecisionAccept, DecisionReject, DecisionNeedsMoreInfo, DecisionEscalate}
}
