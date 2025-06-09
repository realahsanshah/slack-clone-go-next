-- name: JoinWorkspace :one
INSERT INTO workspace_members (
    workspace_id,
    user_id,
    status,
    role
) VALUES (
    $1, $2, $3, $4
) RETURNING id, workspace_id, user_id, status, role, created_at, updated_at, deleted_at;

-- name: LeaveWorkspace :exec
UPDATE workspace_members SET deleted_at = CURRENT_TIMESTAMP WHERE workspace_id = $1 AND user_id = $2;

-- name: GetWorkspaceMembers :many
SELECT * FROM workspace_members WHERE workspace_id = $1 AND deleted_at IS NULL;

-- name: GetUserWorkspaces :many
SELECT * FROM workspace_members WHERE user_id = $1 AND deleted_at IS NULL;

-- name: UpdateMemberStatus :exec
UPDATE workspace_members SET status = $3 WHERE workspace_id = $1 AND user_id = $2;

-- name: UpdateMemberRole :exec
UPDATE workspace_members SET role = $3 WHERE workspace_id = $1 AND user_id = $2;

-- name: GetMemberByWorkspaceIdAndUserId :one
SELECT * FROM workspace_members WHERE workspace_id = $1 AND user_id = $2 AND deleted_at IS NULL;