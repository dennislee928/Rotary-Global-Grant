package dto

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// CreateAPIKeyRequest represents the request body for creating an API key
type CreateAPIKeyRequest struct {
	Name      string   `json:"name" binding:"required,max=255"`
	Scopes    []string `json:"scopes,omitempty"`
	ExpiresIn int      `json:"expiresIn,omitempty"` // Duration in days, 0 = no expiry
}
