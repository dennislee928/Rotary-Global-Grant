package vo

// KPIMetricsVO represents the KPI metrics response
// @Description KPI metrics summary
type KPIMetricsVO struct {
	// Education metrics
	Education EducationKPIVO `json:"education"`
	// System metrics
	System SystemKPIVO `json:"system"`
	// Adoption metrics
	Adoption AdoptionKPIVO `json:"adoption"`
	// Governance metrics
	Governance GovernanceKPIVO `json:"governance"`
}

// EducationKPIVO represents education-related KPIs
// @Description Education KPI metrics
type EducationKPIVO struct {
	// Total workshops delivered
	WorkshopsCount int `json:"workshopsCount" example:"12"`
	// Workshops target
	WorkshopsTarget int `json:"workshopsTarget" example:"12"`
	// Total participants trained
	ParticipantsTrained int `json:"participantsTrained" example:"342"`
	// Participants target
	ParticipantsTarget int `json:"participantsTarget" example:"300"`
	// Average pre/post improvement percentage
	PrePostImprovement float64 `json:"prePostImprovement" example:"25.5"`
	// Improvement target
	ImprovementTarget float64 `json:"improvementTarget" example:"25.0"`
}

// SystemKPIVO represents system-related KPIs
// @Description System KPI metrics
type SystemKPIVO struct {
	// Median report to triage time (minutes)
	MedianReportToTriage float64 `json:"medianReportToTriage" example:"22.5"`
	// Target time (minutes)
	TriageTimeTarget float64 `json:"triageTimeTarget" example:"30.0"`
	// Verified reports ratio (percentage)
	VerifiedRatio float64 `json:"verifiedRatio" example:"65.2"`
	// Verified ratio target
	VerifiedRatioTarget float64 `json:"verifiedRatioTarget" example:"60.0"`
	// Abuse/false report rate (percentage)
	AbuseRate float64 `json:"abuseRate" example:"3.2"`
	// Abuse rate target (max)
	AbuseRateTarget float64 `json:"abuseRateTarget" example:"5.0"`
	// Alert publish latency (minutes)
	AlertPublishLatency float64 `json:"alertPublishLatency" example:"12.0"`
	// Publish latency target
	PublishLatencyTarget float64 `json:"publishLatencyTarget" example:"15.0"`
}

// AdoptionKPIVO represents adoption-related KPIs
// @Description Adoption KPI metrics
type AdoptionKPIVO struct {
	// Partner organizations count
	PartnerOrgs int `json:"partnerOrgs" example:"4"`
	// Partner orgs target
	PartnerOrgsTarget int `json:"partnerOrgsTarget" example:"4"`
	// External deployments/forks count
	ExternalAdoption int `json:"externalAdoption" example:"2"`
	// External adoption target
	ExternalAdoptionTarget int `json:"externalAdoptionTarget" example:"2"`
}

// GovernanceKPIVO represents governance-related KPIs
// @Description Governance KPI metrics
type GovernanceKPIVO struct {
	// Certified triagers count
	CertifiedTriagers int `json:"certifiedTriagers" example:"18"`
	// Certified triagers target
	TriagersTarget int `json:"triagersTarget" example:"15"`
}

// DashboardStatsVO represents public dashboard statistics
// @Description Public dashboard statistics
type DashboardStatsVO struct {
	// Total reports submitted
	TotalReports int `json:"totalReports" example:"156"`
	// Reports this week
	ReportsThisWeek int `json:"reportsThisWeek" example:"12"`
	// Active alerts count
	ActiveAlerts int `json:"activeAlerts" example:"2"`
	// Recent alerts
	RecentAlerts []AlertSummaryVO `json:"recentAlerts"`
	// Report category breakdown
	CategoryBreakdown []CategoryCountVO `json:"categoryBreakdown"`
}

// AlertSummaryVO represents a brief alert summary
// @Description Brief alert summary for dashboard
type AlertSummaryVO struct {
	ID        string `json:"id"`
	Event     string `json:"event"`
	Severity  string `json:"severity"`
	Area      string `json:"area"`
	CreatedAt string `json:"createdAt"`
}

// CategoryCountVO represents a category count pair
// @Description Category count for breakdown
type CategoryCountVO struct {
	Category string `json:"category" example:"scam_phishing"`
	Count    int    `json:"count" example:"45"`
}
