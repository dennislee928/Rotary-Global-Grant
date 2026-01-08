package dto

// CreateReportRequest represents the request body for creating a report
type CreateReportRequest struct {
	Category          string   `json:"category" binding:"required,oneof=suspicious_item suspicious_person harassment_stalking scam_phishing misinformation_panic crowd_disorder infrastructure_hazard other"`
	SeveritySuggested string   `json:"severitySuggested,omitempty" binding:"omitempty,oneof=S0 S1 S2 S3 S4"`
	AreaHint          string   `json:"areaHint" binding:"required,max=500"`
	TimeWindow        string   `json:"timeWindow,omitempty" binding:"max=100"`
	Description       string   `json:"description" binding:"required"`
	Evidence          []string `json:"evidence,omitempty"`
	ReporterContact   string   `json:"reporterContact,omitempty" binding:"max=255"`
}

// ListReportsQuery represents query parameters for listing reports
type ListReportsQuery struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"pageSize,default=20" binding:"min=1,max=100"`
	Status   string `form:"status,omitempty" binding:"omitempty,oneof=submitted under_review triaged escalated closed spam"`
	Category string `form:"category,omitempty" binding:"omitempty,oneof=suspicious_item suspicious_person harassment_stalking scam_phishing misinformation_panic crowd_disorder infrastructure_hazard other"`
	SortBy   string `form:"sortBy,default=createdAt" binding:"oneof=createdAt category status"`
	SortDir  string `form:"sortDir,default=desc" binding:"oneof=asc desc"`
}
