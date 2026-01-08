package vo

import "time"

// ReportVO represents the response for a report
// @Description Report response object
type ReportVO struct {
	// Unique identifier
	ID string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	// Report category
	Category string `json:"category" example:"scam_phishing"`
	// Suggested severity level
	SeveritySuggested string `json:"severitySuggested,omitempty" example:"S2"`
	// Approximate area/location hint
	AreaHint string `json:"areaHint" example:"Near campus entrance"`
	// Time window of the incident
	TimeWindow string `json:"timeWindow,omitempty" example:"2026-01-08 14:00-15:00"`
	// Detailed description
	Description string `json:"description" example:"Suspicious person asking for personal information"`
	// Evidence references (URLs or opaque refs)
	Evidence []string `json:"evidence,omitempty"`
	// Current status
	Status string `json:"status" example:"submitted"`
	// Creation timestamp
	CreatedAt time.Time `json:"createdAt" example:"2026-01-08T14:30:00Z"`
	// Last update timestamp
	UpdatedAt time.Time `json:"updatedAt" example:"2026-01-08T14:30:00Z"`
}

// ReportListVO represents a paginated list of reports
// @Description Paginated report list response
type ReportListVO struct {
	// List of reports
	Data []ReportVO `json:"data"`
	// Pagination metadata
	Pagination PaginationVO `json:"pagination"`
}

// ReportDetailVO represents detailed report with triage history
// @Description Detailed report with triage decisions
type ReportDetailVO struct {
	ReportVO
	// Triage decisions for this report
	TriageDecisions []TriageDecisionVO `json:"triageDecisions,omitempty"`
}
