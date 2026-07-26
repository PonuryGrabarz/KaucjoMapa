-- name: CreateUserNotification :one
INSERT INTO user_notifications (
    id,
    user_id,
    name,
    description,
    seen,
    created_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: ListUserNotificationsByUserId :many
SELECT *
FROM user_notifications
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateUserNotification :one
UPDATE user_notifications
SET
    user_id = $2,
    name = $3,
    description = $4,
    seen = $5
WHERE id = $1
RETURNING *;

-- name: GetNotificationById :one
SELECT * FROM user_notifications WHERE id = $1;

-- name: DeleteUserNotification :exec
DELETE FROM user_notifications
WHERE id = $1;
