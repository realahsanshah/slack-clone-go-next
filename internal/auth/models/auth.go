package models

// RegisterRequest represents the registration request
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=50" example:"John Doe" minLength:"2" maxLength:"50"`
	Email    string `json:"email" binding:"required,email" example:"john.doe@example.com" format:"email"`
	Password string `json:"password" binding:"required,min=6" example:"securepassword123" minLength:"6"`
}

// LoginRequest represents the login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"john.doe@example.com" format:"email"`
	Password string `json:"password" binding:"required" example:"securepassword123"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  User   `json:"user"`
}

// User represents the user data in responses
type User struct {
	ID    int32  `json:"id" example:"1"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john.doe@example.com"`
}
