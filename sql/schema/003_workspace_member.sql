-- +goose Up
-- Create custom enum types for PostgreSQL
CREATE TYPE member_status AS ENUM ('pending', 'accepted', 'rejected');
CREATE TYPE member_role AS ENUM ('admin', 'member');

CREATE TABLE workspace_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workspace_id UUID NOT NULL REFERENCES workspaces(id),
    user_id UUID NOT NULL REFERENCES users(id),
    status member_status NOT NULL DEFAULT 'pending',
    role member_role NOT NULL DEFAULT 'member',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(workspace_id, user_id)
);

-- +goose Down
DROP TABLE workspace_members;
DROP TYPE member_role;
DROP TYPE member_status;