package router

import (
	"github.com/gin-gonic/gin"
)

// PingResponse represents the response for the ping endpoint
type PingResponse struct {
	Message string `json:"message" example:"pong"`
}

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// add base path
	api := router.Group("/api/v1")

	// @Summary Health check endpoint
	// @Description Check if the server is running
	// @Tags health
	// @Accept json
	// @Produce json
	// @Success 200 {object} PingResponse
	// @Router /ping [get]
	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, PingResponse{
			Message: "pong",
		})
	})

	return router
}
