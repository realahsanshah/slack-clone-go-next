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
