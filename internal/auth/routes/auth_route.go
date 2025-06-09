package auth

import (
	controller "slack-clone-go-next/internal/auth/controller"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.RouterGroup) {
	router.POST("/register", controller.Register)
	router.POST("/login", controller.Login)
}
