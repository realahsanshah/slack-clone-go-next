package models

type Workspace struct {
	ID          string `json:"id" example:"1"`
	Name        string `json:"name" example:"John Doe"`
	Username    string `json:"username" example:"john.doe"`
	Logo        string `json:"logo" example:"https://example.com/logo.png" format:"uri"`
	MemberCount int    `json:"member_count" example:"10"`
	UserID      string `json:"user_id" example:"1"`
}

type CreateWorkspaceRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=50" example:"John Doe" minLength:"2" maxLength:"50"`
	Username string `json:"username" binding:"required,min=2,max=50" example:"john.doe" minLength:"2" maxLength:"50"`
	Logo     string `json:"logo" binding:"required,url" example:"https://example.com/logo.png" format:"uri"`
}

type GetWorkspacesRequest struct {
	Page  int `json:"page" binding:"required,min=1" example:"1"`
	Limit int `json:"limit" binding:"required,min=1,max=100" example:"10"`
}
