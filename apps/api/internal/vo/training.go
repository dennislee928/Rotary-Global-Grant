package vo

import "time"

// TrainingEventVO represents the response for a training event
// @Description Training event response object
type TrainingEventVO struct {
	// Unique identifier
	ID string `json:"id" example:"550e8400-e29b-41d4-a716-446655440003"`
	// Event title
	Title string `json:"title" example:"Anti-Fraud Workshop: Recognizing Phishing"`
	// Event date
	EventDate string `json:"eventDate" example:"2026-01-15"`
	// Event location
	Location string `json:"location" example:"Community Center, Room 201"`
	// Target audience
	Audience string `json:"audience,omitempty" example:"Older adults, new residents"`
	// Number of attendees
	AttendanceCount int `json:"attendanceCount" example:"45"`
	// Average pre-test score
	PreAvg *float64 `json:"preAvg,omitempty" example:"62.5"`
	// Average post-test score
	PostAvg *float64 `json:"postAvg,omitempty" example:"85.3"`
	// Improvement percentage
	Improvement *float64 `json:"improvement,omitempty" example:"22.8"`
	// Additional notes
	Notes string `json:"notes,omitempty"`
	// Creation timestamp
	CreatedAt time.Time `json:"createdAt" example:"2026-01-08T10:00:00Z"`
	// Last update timestamp
	UpdatedAt time.Time `json:"updatedAt" example:"2026-01-08T10:00:00Z"`
}

// TrainingEventListVO represents a paginated list of training events
// @Description Paginated training event list response
type TrainingEventListVO struct {
	// List of training events
	Data []TrainingEventVO `json:"data"`
	// Pagination metadata
	Pagination PaginationVO `json:"pagination"`
}

// QuizResultVO represents the response for a quiz result
// @Description Quiz result response object
type QuizResultVO struct {
	// Unique identifier
	ID string `json:"id" example:"550e8400-e29b-41d4-a716-446655440004"`
	// Associated training event ID
	EventID string `json:"eventId,omitempty"`
	// Quiz type (pre/post)
	QuizType string `json:"quizType" example:"post"`
	// Score achieved
	Score float64 `json:"score" example:"85"`
	// Maximum possible score
	MaxScore float64 `json:"maxScore" example:"100"`
	// Score percentage
	Percentage float64 `json:"percentage" example:"85.0"`
	// Timestamp
	CreatedAt time.Time `json:"createdAt" example:"2026-01-15T11:30:00Z"`
}

// TrainingStatsVO represents training statistics summary
// @Description Training statistics summary
type TrainingStatsVO struct {
	// Total events count
	TotalEvents int `json:"totalEvents" example:"12"`
	// Total participants (de-duplicated)
	TotalParticipants int `json:"totalParticipants" example:"342"`
	// Average improvement percentage
	AverageImprovement float64 `json:"averageImprovement" example:"24.5"`
	// Target achievement status
	TargetMet bool `json:"targetMet" example:"true"`
}
