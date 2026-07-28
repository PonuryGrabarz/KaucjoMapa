-- name: CreatePasswordResetChallenge :one
INSERT INTO password_reset_challenges (
    id,
    user_id,
    code_hash,
    expires_at,
    resend_available_at,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;

-- name: GetPasswordResetChallengeByID :one
SELECT * FROM password_reset_challenges WHERE id = $1;

-- name: ConsumeActivePasswordResetChallengesByUserID :exec
UPDATE password_reset_challenges
    SET consumed_at = NOW(), updated_at = NOW()
    WHERE user_id = $1 AND consumed_at IS NULL;

-- name: IncrementPasswordResetChallengeAttempts :one
UPDATE password_reset_challenges
    SET attempts = attempts + 1, updated_at = NOW()
    WHERE id = $1
    RETURNING *;

-- name: UpdatePasswordResetChallengeCode :one
UPDATE password_reset_challenges
    SET code_hash = $2,
        attempts = 0,
        expires_at = $3,
        resend_available_at = $4,
        updated_at = NOW()
    WHERE id = $1
    RETURNING *;

-- name: SetPasswordResetToken :one
UPDATE password_reset_challenges
    SET reset_token_hash = $2,
        reset_token_expires_at = $3,
        updated_at = NOW()
    WHERE id = $1
    RETURNING *;

-- name: ConsumePasswordResetChallenge :exec
UPDATE password_reset_challenges
    SET consumed_at = NOW(), updated_at = NOW()
    WHERE id = $1;
