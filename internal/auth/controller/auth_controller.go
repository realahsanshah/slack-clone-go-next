package auth

import (
	"context"
	"net/http"
	"slack-clone-go-next/internal/auth/models"
	"slack-clone-go-next/internal/database"
	"slack-clone-go-next/middleware"
	"slack-clone-go-next/utils"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Registration details"
// @Success 201 {object} middleware.APIResponse{data=models.AuthResponse} "Registration successful"
// @Failure 400 {object} middleware.APIResponse "Invalid request data"
// @Failure 409 {object} middleware.APIResponse "User with this email already exists"
// @Failure 500 {object} middleware.APIResponse "Internal server error"
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	// Check if user already exists
	_, err := database.DBQueries.GetUserByEmail(context.Background(), req.Email)
	if err == nil {
		middleware.ErrorResponse(c, http.StatusConflict, "User with this email already exists", nil)
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "Failed to process password", err)
		return
	}

	// Create user
	user, err := database.DBQueries.CreateUser(context.Background(), database.CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	})
	if err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}

	response := models.AuthResponse{
		Token: token,
		User: models.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}

	middleware.SuccessResponse(c, response, "Registration successful", http.StatusCreated)
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login credentials"
// @Success 200 {object} middleware.APIResponse{data=models.AuthResponse} "Login successful"
// @Failure 400 {object} middleware.APIResponse "Invalid request data"
// @Failure 401 {object} middleware.APIResponse "Invalid credentials"
// @Failure 500 {object} middleware.APIResponse "Internal server error"
// @Router /auth/login [post]
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	// Get user by email
	user, err := database.DBQueries.GetUserByEmail(context.Background(), req.Email)
	if err != nil {
		middleware.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	// Check password
	if !utils.CheckPassword(req.Password, user.Password) {
		middleware.ErrorResponse(c, http.StatusUnauthorized, "Invalid credentials", nil)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}

	response := models.AuthResponse{
		Token: token,
		User: models.User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}

	middleware.SuccessResponse(c, response, "Login successful", http.StatusOK)
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get the current authenticated user's profile information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} middleware.APIResponse{data=models.User} "Profile retrieved successfully"
// @Failure 401 {object} middleware.APIResponse "Unauthorized - invalid or missing token"
// @Failure 404 {object} middleware.APIResponse "User not found"
// @Router /auth/profile [get]
func GetProfile(c *gin.Context) {
	userEmail, exists := middleware.GetUserEmail(c)
	if !exists {
		middleware.ErrorResponse(c, http.StatusUnauthorized, "User not found in context", nil)
		return
	}

	// Get user from database
	user, err := database.DBQueries.GetUserByEmail(context.Background(), userEmail)
	if err != nil {
		middleware.ErrorResponse(c, http.StatusNotFound, "User not found", err)
		return
	}

	userProfile := models.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	middleware.SuccessResponse(c, userProfile, "Profile retrieved successfully", http.StatusOK)
}
