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

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUserPassword :exec
UPDATE users
    SET password_hash = $2, updated_at = NOW()
    WHERE id = $1;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByGoogleSubject :one
SELECT * FROM users WHERE google_subject = $1;

-- name: UpsertGoogleUserByEmail :one
INSERT INTO users (
    id,
    email,
    username,
    password_hash,
    google_subject
)
VALUES (
    sqlc.arg(id),
    sqlc.arg(email),
    sqlc.arg(username),
    sqlc.arg(password_hash),
    sqlc.arg(google_subject)
)
ON CONFLICT (email) DO UPDATE SET
    google_subject = EXCLUDED.google_subject,
    updated_at = NOW()
WHERE users.google_subject IS NULL
   OR users.google_subject = EXCLUDED.google_subject
RETURNING *;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT * FROM users;

-- name: UpsertAdminUser :exec
INSERT INTO users (
    id,
    email,
    username,
    password_hash,
    reputation,
    is_admin,
    is_banned,
    created_at,
    updated_at
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    0,
    TRUE,
    FALSE,
    NOW(),
    NOW()
)
ON CONFLICT (username) DO UPDATE SET
    email = EXCLUDED.email,
    password_hash = EXCLUDED.password_hash,
    is_admin = TRUE,
    is_banned = FALSE,
    updated_at = NOW();
