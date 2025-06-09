package router

import (
	"net/http"
	auth "slack-clone-go-next/internal/auth/routes"
	"slack-clone-go-next/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// PingResponse represents the response for the ping endpoint
type PingResponse struct {
	Message string `json:"message" example:"pong"`
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Add global response middleware
	router.Use(middleware.ResponseMiddleware())

	// Swagger endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// add base path
	api := router.Group("/api/v1")

	// @Summary Health check endpoint
	// @Description Check if the server is running
	// @Tags health
	// @Accept json
	// @Produce json
	// @Success 200 {object} middleware.APIResponse{data=PingResponse} "Health check successful"
	// @Router /ping [get]
	api.GET("/ping", func(c *gin.Context) {
		middleware.SuccessResponse(c, PingResponse{
			Message: "pong",
		}, "Health check successful", http.StatusOK)
	})

	// add base path for auth routes
	authRoutes := api.Group("/auth")
	auth.RegisterAuthRoutes(authRoutes)

	return router
}
