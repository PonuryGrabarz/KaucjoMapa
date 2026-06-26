-- name: CreateUser :one
INSERT INTO users (id, email, username, password_hash, reputation, is_admin, is_banned, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9
)
RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT * FROM users;

