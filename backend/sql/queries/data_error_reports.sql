-- name: CreateDataErrorReport :one
INSERT INTO data_error_reports (
    id,
    user_id,
    return_point_id,
    categories,
    description,
    suggested_correction,
    status
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

-- name: GetDataErrorReportById :one
SELECT
    der.*,
    u.username,
    rp.name AS return_point_name
FROM data_error_reports AS der
JOIN users AS u ON u.id = der.user_id
JOIN return_points AS rp ON rp.id = der.return_point_id
WHERE der.id = $1;

-- name: ListDataErrorReports :many
SELECT
    der.*,
    u.username,
    rp.name AS return_point_name
FROM data_error_reports AS der
JOIN users AS u ON u.id = der.user_id
JOIN return_points AS rp ON rp.id = der.return_point_id
WHERE (
    sqlc.narg('status')::data_error_report_status IS NULL
    OR der.status = sqlc.narg('status')::data_error_report_status
)
ORDER BY der.created_at DESC;

-- name: ListDataErrorReportsByReturnPoint :many
SELECT
    der.*,
    u.username,
    rp.name AS return_point_name
FROM data_error_reports AS der
JOIN users AS u ON u.id = der.user_id
JOIN return_points AS rp ON rp.id = der.return_point_id
WHERE der.return_point_id = $1
ORDER BY der.created_at DESC;

-- name: ListDataErrorReportsByUserId :many
SELECT
    der.*,
    u.username,
    rp.name AS return_point_name
FROM data_error_reports AS der
JOIN users AS u ON u.id = der.user_id
JOIN return_points AS rp ON rp.id = der.return_point_id
WHERE der.user_id = $1
ORDER BY der.created_at DESC;

-- name: UpdateDataErrorReportStatus :one
UPDATE data_error_reports
SET
    status = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteDataErrorReportById :exec
DELETE FROM data_error_reports
WHERE id = $1;
