-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    id,
    token,
    created_at,
    updated_at,
    user_id,
    expires_at,
    revoked_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
    ) RETURNING *;

-- name: GetUserTokens :many
SELECT * FROM refresh_tokens WHERE user_id = $1;

-- name: GetUserRefreshToken :one
SELECT * FROM refresh_tokens WHERE user_id = $1 AND token = $2 ORDER BY created_at DESC LIMIT 1;

-- name: RevokeToken :exec
UPDATE refresh_tokens
    SET revoked_at = NOW(), updated_at = NOW()
    WHERE token = $1;

-- name: RevokeAllUserTokens :exec
UPDATE refresh_tokens
    SET revoked_at = NOW(), updated_at = NOW()
    WHERE user_id = $1 AND revoked_at IS NULL;
