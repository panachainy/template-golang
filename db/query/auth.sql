-- name: GetAuthByID :one
SELECT * FROM auths
WHERE id = $1;

-- name: GetAuthorByUsername :one
SELECT * FROM auths
WHERE username = $1;

-- name: ListAllAuths :many
SELECT * FROM auths
ORDER BY id;
