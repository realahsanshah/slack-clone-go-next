package workspace

import (
	"database/sql"
	"fmt"
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
// @Failure 409 {object} middleware.APIResponse "Workspace username already exists"
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

	// transaction to create workspace and join workspace
	tx, err := database.DB.BeginTx(c, nil)
	if err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "Failed to begin transaction", err)
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	queries := database.New(tx)
	workspace, err := queries.CreateWorkspace(c, database.CreateWorkspaceParams{
		Name:        req.Name,
		Username:    req.Username,
		Logo:        sql.NullString{String: req.Logo, Valid: req.Logo != ""},
		MemberCount: 1, // Set initial member count to 1 (the creator)
		UserID:      uuid.MustParse(userID.String()),
	})
	if err != nil {
		// Check if it's a duplicate username error
		if err.Error() == "pq: duplicate key value violates unique constraint \"workspaces_username_key\"" {
			middleware.ErrorResponse(c, http.StatusConflict, "Workspace username already exists", err)
			return
		}
		middleware.ErrorResponse(c, http.StatusInternalServerError, "Failed to create workspace", err)
		return
	}

	_, err = queries.JoinWorkspace(c, database.JoinWorkspaceParams{
		WorkspaceID: workspace.ID,
		UserID:      uuid.MustParse(userID.String()),
		Status:      database.MemberStatusAccepted,
		Role:        database.MemberRoleAdmin,
	})
	if err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "Failed to join workspace", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "Failed to commit transaction", err)
		return
	}

	middleware.SuccessResponse(c, models.DatabaseWorkspaceToWorkspace(workspace), "Workspace created successfully", http.StatusCreated)
}

// @Summary Get workspaces
// @Description Gets all workspaces for the authenticated user
// @Tags workspaces
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Page query int false "Page number" default(1)
// @Param Limit query int false "Limit per page" default(10)
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

	page := req.Page
	if page <= 0 {
		page = 1
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	workspaces, err := database.DBQueries.GetUserJoinedWorkspaces(c, database.GetUserJoinedWorkspacesParams{
		UserID: uuid.MustParse(userID.String()),
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	fmt.Println("workspaces", workspaces)
	if err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "Failed to get workspaces", err)
		return
	}

	middleware.SuccessResponse(c, models.DatabaseWorkspacesToWorkspaces(workspaces), "Workspaces fetched successfully", http.StatusOK)
}

// @Summary Get workspace by id
// @Description Gets a workspace by id
// @Tags workspaces
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Workspace ID"
// @Success 200 {object} middleware.APIResponse{data=models.Workspace} "Workspace fetched successfully"
// @Failure 400 {object} middleware.APIResponse "Invalid request data"
// @Failure 401 {object} middleware.APIResponse "Unauthorized"
// @Failure 404 {object} middleware.APIResponse "Workspace not found"
// @Failure 500 {object} middleware.APIResponse "Internal server error"
// @Router /workspaces/{id} [get]
func GetWorkspaceById(c *gin.Context) {
	workspaceID := c.Param("id")
	fmt.Println("workspaceID", workspaceID)
	workspaceIdInUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		middleware.ErrorResponse(c, http.StatusBadRequest, "Invalid workspace ID", err)
		return
	}

	workspace, err := database.DBQueries.GetWorkspaceById(c, workspaceIdInUUID)
	if err != nil {
		middleware.ErrorResponse(c, http.StatusNotFound, "Workspace not found", err)
		return
	}

	middleware.SuccessResponse(c, models.DatabaseWorkspaceToWorkspace(workspace), "Workspace fetched successfully", http.StatusOK)
}

// @Summary Join a workspace
// @Description Joins a workspace for the authenticated user
// @Tags workspaces
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param workspace body models.JoinWorkspaceRequest true "Workspace data"
// @Success 200 {object} middleware.APIResponse "Workspace joined successfully"
// @Failure 400 {object} middleware.APIResponse "Invalid request data"
// @Failure 401 {object} middleware.APIResponse "Unauthorized"
// @Failure 404 {object} middleware.APIResponse "Workspace not found"
// @Failure 500 {object} middleware.APIResponse "Internal server error"
// @Router /workspaces/join [post]
func JoinWorkspace(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		middleware.ErrorResponse(c, http.StatusUnauthorized, "User not found in context", nil)
		return
	}

	var req models.JoinWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	workspaceID := req.WorkspaceID

	workspaceIdInUUID, err := uuid.Parse(workspaceID)
	if err != nil {
		middleware.ErrorResponse(c, http.StatusBadRequest, "Invalid workspace ID", err)
		return
	}

	workspace, err := database.DBQueries.GetWorkspaceById(c, workspaceIdInUUID)
	if err != nil {
		middleware.ErrorResponse(c, http.StatusNotFound, "Workspace not found", err)
		return
	}

	tx, err := database.DB.BeginTx(c, nil)
	if err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "Failed to begin transaction", err)
		return
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	queries := database.New(tx)
	_, err = queries.JoinWorkspace(c, database.JoinWorkspaceParams{
		WorkspaceID: workspace.ID,
		UserID:      uuid.MustParse(userID.String()),
		Status:      database.MemberStatusPending,
		Role:        database.MemberRoleMember,
	})
	if err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "Failed to join workspace", err)
		return
	}

	err = tx.Commit()
	if err != nil {
		middleware.ErrorResponse(c, http.StatusInternalServerError, "Failed to commit transaction", err)
		return
	}

	middleware.SuccessResponse(c, nil, "Workspace joined successfully", http.StatusOK)
}
