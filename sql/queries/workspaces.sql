-- name: CreateWorkspace :one
INSERT INTO workspaces(
    id,
    name,
    username,
    logo,
    member_count,
    user_id
) VALUES (
    $1, $2, $3, $4, $5, $6
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