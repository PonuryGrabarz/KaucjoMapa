-- name: CreatePoint :one
INSERT INTO return_points (
  id,
  user_id,
  name,
  type,
  address,
  latitude,
  longitude,
  opening_hours,
  description,
  current_status,
  confirmations_count,
  rejections_count,
  is_verified,
  created_at,
  updated_at
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8,
  $9,
  $10,
  $11,
  $12,
  $13,
  $14,
  $15
  ) RETURNING *;

-- name: GetAllReturnPoints :many
SELECT * FROM return_points;
