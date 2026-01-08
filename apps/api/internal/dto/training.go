package dto

// CreateTrainingEventRequest represents the request body for creating a training event
type CreateTrainingEventRequest struct {
	Title           string  `json:"title" binding:"required,max=255"`
	EventDate       string  `json:"eventDate" binding:"required"` // Format: YYYY-MM-DD
	Location        string  `json:"location" binding:"required,max=500"`
	Audience        string  `json:"audience,omitempty" binding:"max=255"`
	AttendanceCount int     `json:"attendanceCount,omitempty" binding:"min=0"`
	PreAvg          float64 `json:"preAvg,omitempty" binding:"min=0,max=100"`
	PostAvg         float64 `json:"postAvg,omitempty" binding:"min=0,max=100"`
	Notes           string  `json:"notes,omitempty"`
}

// UpdateTrainingEventRequest represents the request body for updating a training event
type UpdateTrainingEventRequest struct {
	Title           string   `json:"title,omitempty" binding:"omitempty,max=255"`
	EventDate       string   `json:"eventDate,omitempty"` // Format: YYYY-MM-DD
	Location        string   `json:"location,omitempty" binding:"omitempty,max=500"`
	Audience        string   `json:"audience,omitempty" binding:"max=255"`
	AttendanceCount *int     `json:"attendanceCount,omitempty" binding:"omitempty,min=0"`
	PreAvg          *float64 `json:"preAvg,omitempty" binding:"omitempty,min=0,max=100"`
	PostAvg         *float64 `json:"postAvg,omitempty" binding:"omitempty,min=0,max=100"`
	Notes           string   `json:"notes,omitempty"`
}

// RecordQuizResultRequest represents the request body for recording quiz results
type RecordQuizResultRequest struct {
	ParticipantHash string                 `json:"participantHash,omitempty" binding:"max=64"`
	QuizType        string                 `json:"quizType" binding:"required,oneof=pre post"`
	Score           float64                `json:"score" binding:"required,min=0"`
	MaxScore        float64                `json:"maxScore,omitempty" binding:"omitempty,min=1"`
	Answers         map[string]interface{} `json:"answers,omitempty"`
}

// ListTrainingEventsQuery represents query parameters for listing training events
type ListTrainingEventsQuery struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"pageSize,default=20" binding:"min=1,max=100"`
	From     string `form:"from,omitempty"` // Format: YYYY-MM-DD
	To       string `form:"to,omitempty"`   // Format: YYYY-MM-DD
}
