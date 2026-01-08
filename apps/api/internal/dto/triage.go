package dto

// TriageRequest represents the request body for triaging a report
type TriageRequest struct {
	Decision      string `json:"decision" binding:"required,oneof=accept reject needs_more_info escalate"`
	SeverityFinal string `json:"severityFinal" binding:"required,oneof=S0 S1 S2 S3 S4"`
	EvidenceLevel string `json:"evidenceLevel,omitempty" binding:"omitempty,oneof=E0 E1 E2 E3"`
	Rationale     string `json:"rationale,omitempty"`
}

// ListTriageDecisionsQuery represents query parameters for listing triage decisions
type ListTriageDecisionsQuery struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"pageSize,default=20" binding:"min=1,max=100"`
	ReportID string `form:"reportId,omitempty" binding:"omitempty,uuid"`
	Decision string `form:"decision,omitempty" binding:"omitempty,oneof=accept reject needs_more_info escalate"`
}
