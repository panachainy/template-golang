-- name: CreateCockroach :one
INSERT INTO cockroaches (amount)
VALUES ($1)
RETURNING id, amount, created_at;

-- name: GetCockroachByID :one
SELECT id, amount, created_at FROM cockroaches
WHERE id = $1;

-- name: ListCockroaches :many
SELECT id, amount, created_at FROM cockroaches
ORDER BY created_at DESC;

-- name: UpdateCockroach :one
UPDATE cockroaches
SET amount = $2
WHERE id = $1
RETURNING id, amount, created_at;

-- name: DeleteCockroach :exec
DELETE FROM cockroaches
WHERE id = $1;
