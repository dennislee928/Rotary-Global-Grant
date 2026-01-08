package vo

import "time"

// AlertVO represents the response for an alert
// @Description Alert response object
type AlertVO struct {
	// Unique identifier
	ID string `json:"id" example:"550e8400-e29b-41d4-a716-446655440002"`
	// Associated report ID (optional)
	ReportID string `json:"reportId,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	// Alert status
	Status string `json:"status" example:"draft"`
	// Event description
	Event string `json:"event" example:"Suspected phishing campaign targeting campus users"`
	// CAP urgency level
	Urgency string `json:"urgency" example:"Expected"`
	// CAP severity level
	Severity string `json:"severity" example:"Moderate"`
	// CAP certainty level
	Certainty string `json:"certainty" example:"Likely"`
	// Affected area
	Area string `json:"area" example:"University campus and surrounding transit hubs"`
	// Action instructions
	Instruction string `json:"instruction" example:"Do not click suspicious links. Verify sender identity."`
	// Public-facing message
	PublicMessage string `json:"publicMessage,omitempty" example:"Alert: Phishing attempts reported"`
	// Distribution channels
	Channels []string `json:"channels,omitempty" example:"email,sms,web"`
	// CAP XML content
	CAPXML string `json:"capXml,omitempty"`
	// Creation timestamp
	CreatedAt time.Time `json:"createdAt" example:"2026-01-08T15:30:00Z"`
	// Publication timestamp
	PublishedAt *time.Time `json:"publishedAt,omitempty" example:"2026-01-08T16:00:00Z"`
	// Last update timestamp
	UpdatedAt time.Time `json:"updatedAt" example:"2026-01-08T15:30:00Z"`
	// Approver information
	ApprovedBy *UserSummaryVO `json:"approvedBy,omitempty"`
}

// AlertListVO represents a paginated list of alerts
// @Description Paginated alert list response
type AlertListVO struct {
	// List of alerts
	Data []AlertVO `json:"data"`
	// Pagination metadata
	Pagination PaginationVO `json:"pagination"`
}
