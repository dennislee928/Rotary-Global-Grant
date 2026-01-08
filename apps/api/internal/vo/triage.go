package vo

import "time"

// TriageDecisionVO represents the response for a triage decision
// @Description Triage decision response object
type TriageDecisionVO struct {
	// Unique identifier
	ID string `json:"id" example:"550e8400-e29b-41d4-a716-446655440001"`
	// Associated report ID
	ReportID string `json:"reportId" example:"550e8400-e29b-41d4-a716-446655440000"`
	// Decision made
	Decision string `json:"decision" example:"accept"`
	// Final severity assessment
	SeverityFinal string `json:"severityFinal" example:"S2"`
	// Evidence level assessment
	EvidenceLevel string `json:"evidenceLevel,omitempty" example:"E2"`
	// Rationale for the decision
	Rationale string `json:"rationale,omitempty" example:"Clear evidence of phishing attempt"`
	// Decision timestamp
	DecidedAt time.Time `json:"decidedAt" example:"2026-01-08T15:00:00Z"`
	// Decider information (if available)
	DecidedBy *UserSummaryVO `json:"decidedBy,omitempty"`
}

// TriageDecisionListVO represents a paginated list of triage decisions
// @Description Paginated triage decision list response
type TriageDecisionListVO struct {
	// List of triage decisions
	Data []TriageDecisionVO `json:"data"`
	// Pagination metadata
	Pagination PaginationVO `json:"pagination"`
}
