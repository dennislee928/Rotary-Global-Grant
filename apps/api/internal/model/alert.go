package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Alert represents a CAP-ready alert
type Alert struct {
	ID            uuid.UUID   `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ReportID      *uuid.UUID  `gorm:"type:uuid;index"`
	Status        string      `gorm:"size:50;not null;default:'draft'"`
	Event         string      `gorm:"size:255;not null"`
	Urgency       string      `gorm:"size:50;not null"`
	Severity      string      `gorm:"size:50;not null"`
	Certainty     string      `gorm:"size:50;not null"`
	Area          string      `gorm:"size:500;not null"`
	Instruction   string      `gorm:"type:text;not null"`
	PublicMessage string      `gorm:"type:text"`
	CAPXML        string      `gorm:"column:cap_xml;type:text"`
	Channels      StringArray `gorm:"type:jsonb;default:'[]'"`
	ApprovedBy    *uuid.UUID  `gorm:"type:uuid"`
	CreatedAt     time.Time   `gorm:"not null;default:now()"`
	PublishedAt   *time.Time
	UpdatedAt     time.Time `gorm:"not null;default:now()"`

	// Associations
	Report   *Report `gorm:"foreignKey:ReportID"`
	Approver *User   `gorm:"foreignKey:ApprovedBy"`
}

func (Alert) TableName() string {
	return "alerts"
}

func (a *Alert) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.Status == "" {
		a.Status = AlertStatusDraft
	}
	return nil
}

// Alert statuses
const (
	AlertStatusDraft     = "draft"
	AlertStatusApproved  = "approved"
	AlertStatusPublished = "published"
	AlertStatusWithdrawn = "withdrawn"
)

// CAP Urgency values
const (
	UrgencyImmediate = "Immediate"
	UrgencyExpected  = "Expected"
	UrgencyFuture    = "Future"
	UrgencyPast      = "Past"
	UrgencyUnknown   = "Unknown"
)

// CAP Severity values
const (
	CAPSeverityExtreme  = "Extreme"
	CAPSeveritySevere   = "Severe"
	CAPSeverityModerate = "Moderate"
	CAPSeverityMinor    = "Minor"
	CAPSeverityUnknown  = "Unknown"
)

// CAP Certainty values
const (
	CertaintyObserved = "Observed"
	CertaintyLikely   = "Likely"
	CertaintyPossible = "Possible"
	CertaintyUnlikely = "Unlikely"
	CertaintyUnknown  = "Unknown"
)

// ValidAlertStatuses returns all valid alert statuses
func ValidAlertStatuses() []string {
	return []string{AlertStatusDraft, AlertStatusApproved, AlertStatusPublished, AlertStatusWithdrawn}
}

// ValidUrgencies returns all valid CAP urgency values
func ValidUrgencies() []string {
	return []string{UrgencyImmediate, UrgencyExpected, UrgencyFuture, UrgencyPast, UrgencyUnknown}
}

// ValidCAPSeverities returns all valid CAP severity values
func ValidCAPSeverities() []string {
	return []string{CAPSeverityExtreme, CAPSeveritySevere, CAPSeverityModerate, CAPSeverityMinor, CAPSeverityUnknown}
}

// ValidCertainties returns all valid CAP certainty values
func ValidCertainties() []string {
	return []string{CertaintyObserved, CertaintyLikely, CertaintyPossible, CertaintyUnlikely, CertaintyUnknown}
}

// ChannelList is a custom type for handling JSONB channel arrays
type ChannelList []string

func (c ChannelList) Value() (driver.Value, error) {
	if c == nil {
		return "[]", nil
	}
	return json.Marshal(c)
}

func (c *ChannelList) Scan(value interface{}) error {
	if value == nil {
		*c = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan ChannelList")
	}
	return json.Unmarshal(bytes, c)
}
