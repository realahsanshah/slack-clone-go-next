package models

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"Error occurred"`
	Error   string `json:"error,omitempty" example:"Detailed error message"`
	Code    int    `json:"code" example:"400"`
}

// ValidationError represents validation error details
type ValidationError struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"Email is required"`
}
