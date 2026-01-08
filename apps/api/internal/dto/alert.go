package dto

// CreateAlertRequest represents the request body for creating an alert
type CreateAlertRequest struct {
	ReportID      string   `json:"reportId,omitempty" binding:"omitempty,uuid"`
	Event         string   `json:"event" binding:"required,max=255"`
	Urgency       string   `json:"urgency" binding:"required,oneof=Immediate Expected Future Past Unknown"`
	Severity      string   `json:"severity" binding:"required,oneof=Extreme Severe Moderate Minor Unknown"`
	Certainty     string   `json:"certainty" binding:"required,oneof=Observed Likely Possible Unlikely Unknown"`
	Area          string   `json:"area" binding:"required,max=500"`
	Instruction   string   `json:"instruction" binding:"required"`
	PublicMessage string   `json:"publicMessage,omitempty"`
	Channels      []string `json:"channels,omitempty"`
}

// UpdateAlertRequest represents the request body for updating an alert
type UpdateAlertRequest struct {
	Status        string   `json:"status,omitempty" binding:"omitempty,oneof=draft approved published withdrawn"`
	Event         string   `json:"event,omitempty" binding:"omitempty,max=255"`
	Urgency       string   `json:"urgency,omitempty" binding:"omitempty,oneof=Immediate Expected Future Past Unknown"`
	Severity      string   `json:"severity,omitempty" binding:"omitempty,oneof=Extreme Severe Moderate Minor Unknown"`
	Certainty     string   `json:"certainty,omitempty" binding:"omitempty,oneof=Observed Likely Possible Unlikely Unknown"`
	Area          string   `json:"area,omitempty" binding:"omitempty,max=500"`
	Instruction   string   `json:"instruction,omitempty"`
	PublicMessage string   `json:"publicMessage,omitempty"`
	Channels      []string `json:"channels,omitempty"`
}

// ListAlertsQuery represents query parameters for listing alerts
type ListAlertsQuery struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"pageSize,default=20" binding:"min=1,max=100"`
	Status   string `form:"status,omitempty" binding:"omitempty,oneof=draft approved published withdrawn"`
}
