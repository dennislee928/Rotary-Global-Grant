package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Report represents a community incident report
type Report struct {
	ID                 uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Category           string      `gorm:"size:50;not null"`
	SeveritySuggested  string      `gorm:"size:10"`
	AreaHint           string      `gorm:"size:500;not null"`
	TimeWindow         string      `gorm:"size:100"`
	Description        string      `gorm:"type:text;not null"`
	EvidenceRefs       StringArray `gorm:"type:jsonb;default:'[]'"`
	ReporterContactRef string      `gorm:"size:255"`
	Status             string      `gorm:"size:50;not null;default:'submitted'"`
	CreatedAt          time.Time   `gorm:"not null;default:now()"`
	UpdatedAt          time.Time   `gorm:"not null;default:now()"`

	// Associations
	TriageDecisions []TriageDecision `gorm:"foreignKey:ReportID"`
	Alerts          []Alert          `gorm:"foreignKey:ReportID"`
}

func (Report) TableName() string {
	return "reports"
}

func (r *Report) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	if r.Status == "" {
		r.Status = StatusSubmitted
	}
	return nil
}

// Report categories
const (
	CategorySuspiciousItem     = "suspicious_item"
	CategorySuspiciousPerson   = "suspicious_person"
	CategoryHarassmentStalking = "harassment_stalking"
	CategoryScamPhishing       = "scam_phishing"
	CategoryMisinformationPanic = "misinformation_panic"
	CategoryCrowdDisorder      = "crowd_disorder"
	CategoryInfrastructureHazard = "infrastructure_hazard"
	CategoryOther              = "other"
)

// Report statuses
const (
	StatusSubmitted   = "submitted"
	StatusUnderReview = "under_review"
	StatusTriaged     = "triaged"
	StatusEscalated   = "escalated"
	StatusClosed      = "closed"
	StatusSpam        = "spam"
)

// Severity levels
const (
	SeverityS0 = "S0" // Informational
	SeverityS1 = "S1" // Low risk
	SeverityS2 = "S2" // Moderate
	SeverityS3 = "S3" // High
	SeverityS4 = "S4" // Critical
)

// Evidence levels
const (
	EvidenceE0 = "E0" // None
	EvidenceE1 = "E1" // Weak
	EvidenceE2 = "E2" // Moderate
	EvidenceE3 = "E3" // Strong
)

// ValidCategories returns all valid report categories
func ValidCategories() []string {
	return []string{
		CategorySuspiciousItem,
		CategorySuspiciousPerson,
		CategoryHarassmentStalking,
		CategoryScamPhishing,
		CategoryMisinformationPanic,
		CategoryCrowdDisorder,
		CategoryInfrastructureHazard,
		CategoryOther,
	}
}

// ValidSeverities returns all valid severity levels
func ValidSeverities() []string {
	return []string{SeverityS0, SeverityS1, SeverityS2, SeverityS3, SeverityS4}
}

// ValidEvidenceLevels returns all valid evidence levels
func ValidEvidenceLevels() []string {
	return []string{EvidenceE0, EvidenceE1, EvidenceE2, EvidenceE3}
}

// StringArray is a custom type for handling JSONB string arrays
type StringArray []string

func (a StringArray) Value() (driver.Value, error) {
	if a == nil {
		return "[]", nil
	}
	return json.Marshal(a)
}

func (a *StringArray) Scan(value interface{}) error {
	if value == nil {
		*a = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan StringArray")
	}
	return json.Unmarshal(bytes, a)
}
