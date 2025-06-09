-- name: CreateWorkspace :one
INSERT INTO workspaces(
    name,
    username,
    logo,
    member_count,
    user_id
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetWorkspaceByUsername :one
SELECT * FROM workspaces WHERE username = $1;

-- name: GetWorkspaceById :one
SELECT * FROM workspaces WHERE id = $1;

-- name: GetWorkspacesByUserId :many
SELECT * FROM workspaces 
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2
OFFSET $3;

-- name: GetUserJoinedWorkspaces :many
SELECT w.*, wm.status, wm.role, wm.created_at, wm.updated_at, wm.deleted_at FROM workspaces w
LEFT JOIN workspace_members wm ON w.id = wm.workspace_id
WHERE wm.user_id = $1 AND wm.status IN ('accepted', 'pending')
AND wm.deleted_at IS NULL
ORDER BY w.created_at DESC
LIMIT $2
OFFSET $3;

-- name: JoinWorkspace :one
INSERT INTO workspace_members (
    workspace_id,
    user_id,
    status,
    role
) VALUES (
    $1, $2, $3, $4
) RETURNING *;