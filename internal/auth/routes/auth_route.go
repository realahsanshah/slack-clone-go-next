package auth

import (
	controller "slack-clone-go-next/internal/auth/controller"
	"slack-clone-go-next/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.RouterGroup) {
	router.POST("/register", controller.Register)
	router.POST("/login", controller.Login)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/profile", controller.GetProfile)
	}
}
