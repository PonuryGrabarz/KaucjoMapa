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
  glass,
  refund_methods
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
  $18,
  $19
  ) RETURNING *;

-- name: GetAllReturnPoints :many
SELECT * FROM return_points;

-- name: ListReturnPoints :many
SELECT return_points.*,
    CAST(CASE
        WHEN sqlc.narg('latitude')::DOUBLE PRECISION IS NULL THEN 0.0
        ELSE earth_distance(
            ll_to_earth(return_points.latitude, return_points.longitude),
            ll_to_earth(
                sqlc.narg('latitude')::DOUBLE PRECISION,
                sqlc.narg('longitude')::DOUBLE PRECISION
            )
        ) / 1000.0
    END AS DOUBLE PRECISION) AS distance_km
FROM return_points
WHERE (
    sqlc.narg('query')::TEXT IS NULL
    OR name ILIKE '%' || sqlc.narg('query')::TEXT || '%'
    OR COALESCE(address, '') ILIKE '%' || sqlc.narg('query')::TEXT || '%'
)
AND (
    sqlc.narg('point_type')::point_type IS NULL
    OR type = sqlc.narg('point_type')::point_type
)
AND (
    sqlc.narg('point_status')::point_status IS NULL
    OR current_status = sqlc.narg('point_status')::point_status
)
AND (
    sqlc.narg('user_id')::UUID IS NULL
    OR return_points.user_id = sqlc.narg('user_id')::UUID
)
AND (
    sqlc.narg('pet')::BOOLEAN IS NULL
    OR pet = sqlc.narg('pet')::BOOLEAN
)
AND (
    sqlc.narg('can')::BOOLEAN IS NULL
    OR can = sqlc.narg('can')::BOOLEAN
)
AND (
    sqlc.narg('glass')::BOOLEAN IS NULL
    OR glass = sqlc.narg('glass')::BOOLEAN
)
AND (
    sqlc.narg('refund_methods')::refund_method[] IS NULL
    OR return_points.refund_methods && sqlc.narg('refund_methods')::refund_method[]
)
AND (
    sqlc.narg('min_distance_km')::DOUBLE PRECISION IS NULL
    OR earth_distance(
        ll_to_earth(return_points.latitude, return_points.longitude),
        ll_to_earth(
            sqlc.narg('latitude')::DOUBLE PRECISION,
            sqlc.narg('longitude')::DOUBLE PRECISION
        )
    ) >= sqlc.narg('min_distance_km')::DOUBLE PRECISION * 1000.0
)
AND (
    sqlc.narg('max_distance_km')::DOUBLE PRECISION IS NULL
    OR (
        earth_box(
            ll_to_earth(
                sqlc.narg('latitude')::DOUBLE PRECISION,
                sqlc.narg('longitude')::DOUBLE PRECISION
            ),
            sqlc.narg('max_distance_km')::DOUBLE PRECISION * 1000.0
        ) @> ll_to_earth(return_points.latitude, return_points.longitude)
        AND earth_distance(
            ll_to_earth(return_points.latitude, return_points.longitude),
            ll_to_earth(
                sqlc.narg('latitude')::DOUBLE PRECISION,
                sqlc.narg('longitude')::DOUBLE PRECISION
            )
        ) <= sqlc.narg('max_distance_km')::DOUBLE PRECISION * 1000.0
    )
)
ORDER BY distance_km ASC NULLS LAST, id
LIMIT sqlc.arg('limit')::INT
OFFSET sqlc.arg('offset')::INT;

-- name: CountReturnPoints :one
SELECT COUNT(*)
FROM return_points
WHERE (
    sqlc.narg('query')::TEXT IS NULL
    OR name ILIKE '%' || sqlc.narg('query')::TEXT || '%'
    OR COALESCE(address, '') ILIKE '%' || sqlc.narg('query')::TEXT || '%'
)
AND (
    sqlc.narg('point_type')::point_type IS NULL
    OR type = sqlc.narg('point_type')::point_type
)
AND (
    sqlc.narg('point_status')::point_status IS NULL
    OR current_status = sqlc.narg('point_status')::point_status
)
AND (
    sqlc.narg('user_id')::UUID IS NULL
    OR return_points.user_id = sqlc.narg('user_id')::UUID
)
AND (
    sqlc.narg('pet')::BOOLEAN IS NULL
    OR pet = sqlc.narg('pet')::BOOLEAN
)
AND (
    sqlc.narg('can')::BOOLEAN IS NULL
    OR can = sqlc.narg('can')::BOOLEAN
)
AND (
    sqlc.narg('glass')::BOOLEAN IS NULL
    OR glass = sqlc.narg('glass')::BOOLEAN
)
AND (
    sqlc.narg('refund_methods')::refund_method[] IS NULL
    OR return_points.refund_methods && sqlc.narg('refund_methods')::refund_method[]
)
AND (
    sqlc.narg('min_distance_km')::DOUBLE PRECISION IS NULL
    OR earth_distance(
        ll_to_earth(return_points.latitude, return_points.longitude),
        ll_to_earth(
            sqlc.narg('latitude')::DOUBLE PRECISION,
            sqlc.narg('longitude')::DOUBLE PRECISION
        )
    ) >= sqlc.narg('min_distance_km')::DOUBLE PRECISION * 1000.0
)
AND (
    sqlc.narg('max_distance_km')::DOUBLE PRECISION IS NULL
    OR (
        earth_box(
            ll_to_earth(
                sqlc.narg('latitude')::DOUBLE PRECISION,
                sqlc.narg('longitude')::DOUBLE PRECISION
            ),
            sqlc.narg('max_distance_km')::DOUBLE PRECISION * 1000.0
        ) @> ll_to_earth(return_points.latitude, return_points.longitude)
        AND earth_distance(
            ll_to_earth(return_points.latitude, return_points.longitude),
            ll_to_earth(
                sqlc.narg('latitude')::DOUBLE PRECISION,
                sqlc.narg('longitude')::DOUBLE PRECISION
            )
        ) <= sqlc.narg('max_distance_km')::DOUBLE PRECISION * 1000.0
    )
);

-- name: GetPointByDistance :one
SELECT id, name, address FROM return_points
WHERE earth_distance(
    ll_to_earth(latitude, longitude),
    ll_to_earth($1, $2)
) < 25
LIMIT 1;

-- name: GetPointById :one
SELECT * FROM return_points WHERE id = $1;

-- name: DeletePointById :exec
DELETE FROM return_points WHERE id = $1;

-- name: UpdatePointById :one
UPDATE return_points AS rp
SET
    name = COALESCE(
        sqlc.narg('name')::VARCHAR(255),
        rp.name
    ),
    type = COALESCE(
        sqlc.narg('type')::point_type,
        rp.type
    ),
    address = COALESCE(
        sqlc.narg('address')::TEXT,
        rp.address
    ),
    latitude = COALESCE(
        sqlc.narg('latitude')::DOUBLE PRECISION,
        rp.latitude
    ),
    longitude = COALESCE(
        sqlc.narg('longitude')::DOUBLE PRECISION,
        rp.longitude
    ),
    opening_hours = COALESCE(
        sqlc.narg('opening_hours')::JSONB,
        rp.opening_hours
    ),
    description = COALESCE(
        sqlc.narg('description')::TEXT,
        rp.description
    ),
    pet = COALESCE(
        sqlc.narg('pet')::BOOLEAN,
        rp.pet
    ),
    can = COALESCE(
        sqlc.narg('can')::BOOLEAN,
        rp.can
    ),
    refund_methods = COALESCE(
      sqlc.narg('refund_methods')::refund_method[],
      rp.refund_methods
    ),
    glass = COALESCE(
        sqlc.narg('glass')::BOOLEAN,
        rp.glass
    ),
    is_verified = false,
    updated_at = NOW()
WHERE rp.id = sqlc.arg('id')::UUID
RETURNING rp.*;
