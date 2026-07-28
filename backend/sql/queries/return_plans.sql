-- name: CreateReturnPlan :one
INSERT INTO return_plans (
    id,
    user_id,
    status,
    target_pet,
    target_can,
    target_glass,
    rate_pet_grosze,
    rate_can_grosze,
    rate_glass_grosze,
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
    $11
)
RETURNING *;

-- name: GetReturnPlanByID :one
SELECT *
FROM return_plans
WHERE id = $1;

-- name: GetReturnPlanByIDAndUserID :one
SELECT *
FROM return_plans
WHERE id = $1
  AND user_id = $2;

-- name: GetDraftReturnPlanByUserID :one
SELECT *
FROM return_plans
WHERE user_id = $1
  AND status = 'draft';

-- name: ListReturnPlansByUserID :many
SELECT *
FROM return_plans
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateReturnPlanTargets :one
UPDATE return_plans
SET target_pet = $2,
    target_can = $3,
    target_glass = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateReturnPlanStatus :one
UPDATE return_plans
SET status = $2,
    confirmed_at = CASE
        WHEN $2::return_plan_status = 'confirmed' THEN COALESCE(confirmed_at, NOW())
        ELSE confirmed_at
    END,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: ConfirmReturnPlan :one
UPDATE return_plans
SET status = 'confirmed',
    confirmed_at = NOW(),
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteReturnPlanByIDAndUserID :exec
DELETE FROM return_plans
WHERE id = $1
  AND user_id = $2;

-- name: CreateReturnPlanPoint :one
INSERT INTO return_plan_points (
    id,
    return_plan_id,
    return_point_id,
    point_name,
    point_address,
    pet_qty,
    can_qty,
    glass_qty,
    status,
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
    $11
)
RETURNING *;

-- name: GetReturnPlanPointByID :one
SELECT *
FROM return_plan_points
WHERE id = $1;

-- name: GetReturnPlanPointByPlanAndPoint :one
SELECT *
FROM return_plan_points
WHERE return_plan_id = $1
  AND return_point_id = $2;

-- name: ListReturnPlanPointsByPlanID :many
SELECT *
FROM return_plan_points
WHERE return_plan_id = $1
ORDER BY created_at ASC;

-- name: UpdateReturnPlanPointQuantities :one
UPDATE return_plan_points
SET pet_qty = $2,
    can_qty = $3,
    glass_qty = $4,
    updated_at = NOW()
WHERE id = $1
  AND status = 'pending'
RETURNING *;

-- name: ConfirmReturnPlanPoint :one
UPDATE return_plan_points
SET status = 'confirmed',
    confirmed_at = NOW(),
    updated_at = NOW()
WHERE id = $1
  AND status = 'pending'
RETURNING *;

-- name: ConfirmAllPendingReturnPlanPoints :many
UPDATE return_plan_points
SET status = 'confirmed',
    confirmed_at = NOW(),
    updated_at = NOW()
WHERE return_plan_id = $1
  AND status = 'pending'
RETURNING *;

-- name: DeleteReturnPlanPoint :exec
DELETE FROM return_plan_points
WHERE id = $1
  AND status = 'pending';

-- name: SumReturnPlanPointQuantities :one
SELECT
    COALESCE(SUM(pet_qty), 0)::INT AS pet_qty,
    COALESCE(SUM(can_qty), 0)::INT AS can_qty,
    COALESCE(SUM(glass_qty), 0)::INT AS glass_qty
FROM return_plan_points
WHERE return_plan_id = $1;

-- name: CountReturnPlanPointsByStatus :one
SELECT
    COUNT(*) FILTER (WHERE status = 'pending')::INT AS pending_count,
    COUNT(*) FILTER (WHERE status = 'confirmed')::INT AS confirmed_count,
    COUNT(*)::INT AS total_count
FROM return_plan_points
WHERE return_plan_id = $1;
