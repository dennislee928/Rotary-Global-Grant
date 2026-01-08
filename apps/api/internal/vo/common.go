package vo

import "time"

// PaginationVO represents pagination metadata
// @Description Pagination metadata
type PaginationVO struct {
	// Current page number
	Page int `json:"page" example:"1"`
	// Items per page
	PageSize int `json:"pageSize" example:"20"`
	// Total number of items
	Total int64 `json:"total" example:"156"`
	// Total number of pages
	TotalPages int `json:"totalPages" example:"8"`
}

// UserSummaryVO represents a brief user summary
// @Description Brief user information
type UserSummaryVO struct {
	// User ID
	ID string `json:"id" example:"550e8400-e29b-41d4-a716-446655440005"`
	// Display name
	DisplayName string `json:"displayName,omitempty" example:"John Doe"`
	// User role
	Role string `json:"role" example:"triager"`
}

// ErrorVO represents an error response
// @Description Error response
type ErrorVO struct {
	// Error code
	Code string `json:"code" example:"VALIDATION_ERROR"`
	// Error message
	Message string `json:"message" example:"Invalid request parameters"`
	// Field-specific errors (for validation)
	Details map[string]string `json:"details,omitempty"`
}

// SuccessVO represents a success response without data
// @Description Success response
type SuccessVO struct {
	// Success message
	Message string `json:"message" example:"Operation completed successfully"`
}

// HealthVO represents health check response
// @Description Health check response
type HealthVO struct {
	// Service status
	Status string `json:"status" example:"ok"`
	// Current timestamp
	Timestamp time.Time `json:"ts" example:"2026-01-08T14:30:00Z"`
	// Version info
	Version string `json:"version,omitempty" example:"0.1.0"`
	// Database status
	Database string `json:"database,omitempty" example:"connected"`
	// Redis status
	Redis string `json:"redis,omitempty" example:"connected"`
}

// TokenVO represents authentication token response
// @Description Authentication token response
type TokenVO struct {
	// Access token
	AccessToken string `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIs..."`
	// Token type
	TokenType string `json:"tokenType" example:"Bearer"`
	// Expiration time in seconds
	ExpiresIn int `json:"expiresIn" example:"3600"`
}

// APIKeyVO represents API key response
// @Description API key response (key only shown once on creation)
type APIKeyVO struct {
	// API key ID
	ID string `json:"id" example:"550e8400-e29b-41d4-a716-446655440006"`
	// API key name
	Name string `json:"name" example:"CI/CD Pipeline"`
	// API key (only shown on creation)
	Key string `json:"key,omitempty" example:"hive_ak_xxxxxxxxxxxxx"`
	// Scopes
	Scopes []string `json:"scopes" example:"read,write"`
	// Expiration date
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
	// Creation date
	CreatedAt time.Time `json:"createdAt" example:"2026-01-08T10:00:00Z"`
	// Last used date
	LastUsedAt *time.Time `json:"lastUsedAt,omitempty"`
	// Active status
	IsActive bool `json:"isActive" example:"true"`
}

// AuditLogVO represents an audit log entry
// @Description Audit log entry
type AuditLogVO struct {
	// Log entry ID
	ID string `json:"id"`
	// Actor information
	Actor *UserSummaryVO `json:"actor,omitempty"`
	// Actor IP address
	ActorIP string `json:"actorIp,omitempty"`
	// Action performed
	Action string `json:"action" example:"create"`
	// Object type
	ObjectType string `json:"objectType" example:"report"`
	// Object ID
	ObjectID string `json:"objectId,omitempty"`
	// Changes diff
	Diff map[string]interface{} `json:"diff,omitempty"`
	// Timestamp
	Timestamp time.Time `json:"ts" example:"2026-01-08T14:30:00Z"`
}

// AuditLogListVO represents a paginated list of audit logs
// @Description Paginated audit log list
type AuditLogListVO struct {
	// List of audit logs
	Data []AuditLogVO `json:"data"`
	// Pagination metadata
	Pagination PaginationVO `json:"pagination"`
}
