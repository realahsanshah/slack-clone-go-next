package workspace

import (
	controller "slack-clone-go-next/internal/workspace/controller"
	"slack-clone-go-next/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterWorkspaceRoutes(router *gin.RouterGroup) {
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/", controller.CreateWorkspace)
		protected.GET("/", controller.GetWorkspaces)
	}
}
