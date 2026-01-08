package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TrainingEvent represents a training/workshop event
type TrainingEvent struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Title           string     `gorm:"size:255;not null"`
	EventDate       time.Time  `gorm:"type:date;not null"`
	Location        string     `gorm:"size:500;not null"`
	Audience        string     `gorm:"size:255"`
	AttendanceCount int        `gorm:"default:0"`
	PreAvg          *float64   `gorm:"type:numeric(5,2)"`
	PostAvg         *float64   `gorm:"type:numeric(5,2)"`
	Notes           string     `gorm:"type:text"`
	CreatedBy       *uuid.UUID `gorm:"type:uuid"`
	CreatedAt       time.Time  `gorm:"not null;default:now()"`
	UpdatedAt       time.Time  `gorm:"not null;default:now()"`

	// Associations
	Creator      *User                 `gorm:"foreignKey:CreatedBy"`
	Participants []TrainingParticipant `gorm:"foreignKey:EventID"`
	QuizResults  []QuizResult          `gorm:"foreignKey:EventID"`
}

func (TrainingEvent) TableName() string {
	return "training_events"
}

func (t *TrainingEvent) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// TrainingParticipant represents a participant in a training event
type TrainingParticipant struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	EventID         uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_event_participant"`
	ParticipantHash string    `gorm:"size:64;not null;uniqueIndex:idx_event_participant"`
	PreScore        *float64  `gorm:"type:numeric(5,2)"`
	PostScore       *float64  `gorm:"type:numeric(5,2)"`
	CreatedAt       time.Time `gorm:"not null;default:now()"`

	// Associations
	Event TrainingEvent `gorm:"foreignKey:EventID"`
}

func (TrainingParticipant) TableName() string {
	return "training_participants"
}

func (tp *TrainingParticipant) BeforeCreate(tx *gorm.DB) error {
	if tp.ID == uuid.Nil {
		tp.ID = uuid.New()
	}
	return nil
}

// QuizResult represents a quiz result
type QuizResult struct {
	ID              uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	EventID         *uuid.UUID `gorm:"type:uuid;index"`
	ParticipantHash string     `gorm:"size:64"`
	QuizType        string     `gorm:"size:50;not null"`
	Score           float64    `gorm:"type:numeric(5,2);not null"`
	MaxScore        float64    `gorm:"type:numeric(5,2);not null;default:100"`
	Answers         JSONMap    `gorm:"type:jsonb"`
	CreatedAt       time.Time  `gorm:"not null;default:now()"`

	// Associations
	Event *TrainingEvent `gorm:"foreignKey:EventID"`
}

func (QuizResult) TableName() string {
	return "quiz_results"
}

func (q *QuizResult) BeforeCreate(tx *gorm.DB) error {
	if q.ID == uuid.Nil {
		q.ID = uuid.New()
	}
	return nil
}

// Quiz types
const (
	QuizTypePre  = "pre"
	QuizTypePost = "post"
)

// ValidQuizTypes returns all valid quiz types
func ValidQuizTypes() []string {
	return []string{QuizTypePre, QuizTypePost}
}

// JSONMap is a custom type for handling JSONB maps
type JSONMap map[string]interface{}

func (m JSONMap) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	return json.Marshal(m)
}

func (m *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*m = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan JSONMap")
	}
	return json.Unmarshal(bytes, m)
}
