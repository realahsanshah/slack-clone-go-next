package models

import (
	"slack-clone-go-next/internal/database"
)

type Workspace struct {
	ID          string `json:"id" example:"1"`
	Name        string `json:"name" example:"John Doe"`
	Username    string `json:"username" example:"john.doe"`
	Logo        string `json:"logo" example:"https://example.com/logo.png" format:"uri"`
	MemberCount int    `json:"member_count" example:"10"`
	UserID      string `json:"user_id" example:"1"`
	Status      string `json:"status" example:"pending"`
}

type WorkspaceMember struct {
	ID          string `json:"id" example:"1"`
	WorkspaceID string `json:"workspace_id" example:"1"`
	UserID      string `json:"user_id" example:"1"`
	Status      string `json:"status" example:"pending"`
	Role        string `json:"role" example:"member"`
}

type CreateWorkspaceRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=50" example:"John Doe" minLength:"2" maxLength:"50"`
	Username string `json:"username" binding:"required,min=2,max=50" example:"john.doe" minLength:"2" maxLength:"50"`
	Logo     string `json:"logo" binding:"required,url" example:"https://example.com/logo.png" format:"uri"`
}

type JoinWorkspaceRequest struct {
	WorkspaceID string `json:"workspace_id" example:"1"`
}

type GetWorkspacesRequest struct {
	Page  int `json:"page" binding:"required,min=1" example:"1" format:"int32"`
	Limit int `json:"limit" binding:"required,min=1,max=100" example:"10" format:"int32"`
}

type GetWorkspaceMembersRequest struct {
	WorkspaceID string `json:"workspace_id" example:"1"`
}

// DatabaseWorkspaceToWorkspace converts database.Workspace to models.Workspace
func DatabaseWorkspaceToWorkspace(dbWorkspace database.Workspace) Workspace {
	return Workspace{
		ID:          dbWorkspace.ID.String(),
		Name:        dbWorkspace.Name,
		Username:    dbWorkspace.Username,
		Logo:        dbWorkspace.Logo.String,
		MemberCount: int(dbWorkspace.MemberCount),
		UserID:      dbWorkspace.UserID.String(),
	}
}

// DatabaseWorkspacesToWorkspaces converts []database.Workspace to []models.Workspace
func DatabaseWorkspacesToWorkspaces(dbWorkspaces []database.Workspace) []Workspace {
	workspaces := make([]Workspace, len(dbWorkspaces))
	for i, dbWorkspace := range dbWorkspaces {
		workspaces[i] = DatabaseWorkspaceToWorkspace(dbWorkspace)
	}
	return workspaces
}

// DatabaseWorkspaceMemberToWorkspaceMember converts database.WorkspaceMember to models.WorkspaceMember
func DatabaseWorkspaceMemberToWorkspaceMember(dbWorkspaceMember database.WorkspaceMember) WorkspaceMember {
	return WorkspaceMember{
		ID:          dbWorkspaceMember.ID.String(),
		WorkspaceID: dbWorkspaceMember.WorkspaceID.String(),
		UserID:      dbWorkspaceMember.UserID.String(),
		Status:      string(dbWorkspaceMember.Status),
		Role:        string(dbWorkspaceMember.Role),
	}
}

// DatabaseWorkspaceMembersToWorkspaceMembers converts []database.WorkspaceMember to []models.WorkspaceMember
func DatabaseWorkspaceMembersToWorkspaceMembers(dbWorkspaceMembers []database.WorkspaceMember) []WorkspaceMember {
	workspaceMembers := make([]WorkspaceMember, len(dbWorkspaceMembers))
	for i, dbWorkspaceMember := range dbWorkspaceMembers {
		workspaceMembers[i] = DatabaseWorkspaceMemberToWorkspaceMember(dbWorkspaceMember)
	}
	return workspaceMembers
}
