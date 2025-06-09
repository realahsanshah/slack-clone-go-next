package workspace

import (
	"database/sql"
	"net/http"

	"slack-clone-go-next/internal/database"
	"slack-clone-go-next/internal/workspace/models"
	"slack-clone-go-next/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary Create a new workspace
// @Description Creates a new workspace for the authenticated user
// @Tags workspaces
// @Accept json
// @Security BearerAuth
// @Produce json
// @Param workspace body models.CreateWorkspaceRequest true "Workspace data"
// @Success 201 {object} middleware.APIResponse{data=models.Workspace} "Workspace created successfully"
// @Failure 400 {object} middleware.APIResponse "Invalid request data"
// @Failure 401 {object} middleware.APIResponse "Unauthorized"
// @Failure 500 {object} middleware.APIResponse "Internal server error"
// @Router /workspaces [post]
func CreateWorkspace(c *gin.Context) {
	var req models.CreateWorkspaceRequest
	userID, exists := middleware.GetUserID(c)
	if !exists {
		middleware.ErrorResponse(c, http.StatusUnauthorized, "User not found in context", nil)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	workspace, err := database.DBQueries.CreateWorkspace(c, database.CreateWorkspaceParams{
		Name:     req.Name,
		Username: req.Username,
		Logo:     sql.NullString{String: req.Logo, Valid: true},
		UserID:   uuid.MustParse(userID.String()),
	})
	if err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "Failed to create workspace", err)
		return
	}

	middleware.SuccessResponse(c, workspace, "Workspace created successfully", http.StatusCreated)
}

// @Summary Get workspaces
// @Description Gets all workspaces for the authenticated user
// @Tags workspaces
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Limit per page" default(10)
// @Success 200 {object} middleware.APIResponse{data=[]models.Workspace} "Workspaces fetched successfully"
// @Failure 400 {object} middleware.APIResponse "Invalid request data"
// @Failure 401 {object} middleware.APIResponse "Unauthorized"
// @Failure 500 {object} middleware.APIResponse "Internal server error"
// @Router /workspaces [get]
func GetWorkspaces(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		middleware.ErrorResponse(c, http.StatusUnauthorized, "User not found in context", nil)
		return
	}

	var req models.GetWorkspacesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		middleware.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	offset := (req.Page - 1) * limit

	workspaces, err := database.DBQueries.GetWorkspacesByUserId(c, database.GetWorkspacesByUserIdParams{
		UserID: uuid.MustParse(userID.String()),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "Failed to get workspaces", err)
		return
	}

	middleware.SuccessResponse(c, workspaces, "Workspaces fetched successfully", http.StatusOK)
}
