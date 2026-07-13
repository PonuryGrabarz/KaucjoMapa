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
  updated_at,
  pet,
  can,
  glass
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
  $15,
  $16,
  $17,
  $18
  ) RETURNING *;

-- name: GetAllReturnPoints :many
SELECT * FROM return_points;

-- name: GetPointByDistance :one
SELECT id, name, address FROM return_points
WHERE earth_distance(
    ll_to_earth(latitude, longitude),
    ll_to_earth($1, $2)
) < 25
LIMIT 1;
