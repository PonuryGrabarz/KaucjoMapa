-- name: CreateBookmark :one
INSERT INTO bookmarks (
    id,
    user_id,
    return_point_id,
    created_at
) VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetBookmarkById :one
SELECT *
FROM bookmarks
WHERE id = $1;

-- name: GetBookmarkByUserAndReturnPoint :one
SELECT *
FROM bookmarks
WHERE user_id = $1
  AND return_point_id = $2;

-- name: ListBookmarksByUserId :many
SELECT *
FROM bookmarks
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: ListBookmarksWithReturnPointsByUserId :many
SELECT
    b.*,
    rp.name,
    rp.type,
    rp.address,
    rp.latitude,
    rp.longitude,
    rp.pet,
    rp.can,
    rp.glass,
    rp.opening_hours,
    rp.description,
    rp.current_status,
    rp.confirmations_count,
    rp.rejections_count,
    rp.is_verified,
    rp.refund_methods
FROM bookmarks AS b
JOIN return_points AS rp ON rp.id = b.return_point_id
WHERE b.user_id = $1
ORDER BY b.created_at DESC;

-- name: UpdateBookmark :one
UPDATE bookmarks
SET return_point_id = $2
WHERE id = $1
  AND user_id = $3
RETURNING *;

-- name: DeleteBookmarkById :exec
DELETE FROM bookmarks
WHERE id = $1
  AND user_id = $2;

-- name: DeleteBookmarkByUserAndReturnPoint :exec
DELETE FROM bookmarks
WHERE user_id = $1
  AND return_point_id = $2;
