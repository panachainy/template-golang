-- name: GetAuthByID :one
SELECT * FROM auths
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetAuthByUsername :one
SELECT * FROM auths
WHERE username = $1 AND deleted_at IS NULL;

-- name: GetAuthByEmail :one
SELECT * FROM auths
WHERE email = $1 AND deleted_at IS NULL;

-- name: CreateAuth :one
INSERT INTO auths (username, password, email, role, active)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateAuth :one
UPDATE auths
SET username = $2, password = $3, email = $4, role = $5, active = $6, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteAuth :exec
UPDATE auths
SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListAllAuths :many
SELECT * FROM auths
WHERE deleted_at IS NULL
ORDER BY created_at DESC;

-- Auth Methods queries
-- name: CreateAuthMethod :one
INSERT INTO auth_methods (
    auth_id, provider, provider_id, email, user_id, name, first_name, last_name,
    nick_name, description, avatar_url, location, access_token, refresh_token,
    id_token, expires_at, access_token_secret
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
RETURNING *;

-- name: GetAuthMethodByProviderAndID :one
SELECT * FROM auth_methods
WHERE provider = $1 AND provider_id = $2 AND deleted_at IS NULL;

-- name: GetAuthMethodsByAuthID :many
SELECT * FROM auth_methods
WHERE auth_id = $1 AND deleted_at IS NULL;

-- name: UpdateAuthMethod :one
UPDATE auth_methods
SET access_token = $3, refresh_token = $4, id_token = $5, expires_at = $6, updated_at = CURRENT_TIMESTAMP
WHERE auth_id = $1 AND provider = $2 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteAuthMethod :exec
UPDATE auth_methods
SET deleted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = $1 AND deleted_at IS NULL;
