-- name: CreateCommunityReport :one
INSERT INTO community_reports (
    id,
    user_id,
    return_point_id,
    category,
    description
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetCommunityReportById :one
SELECT *
FROM community_reports
WHERE id = $1;

-- name: ListCommunityReportsByReturnPoint :many
SELECT
    cr.*,
    u.username
FROM community_reports AS cr
JOIN users AS u ON u.id = cr.user_id
WHERE cr.return_point_id = $1
ORDER BY cr.created_at DESC;

-- name: DeleteCommunityReportById :exec
DELETE FROM community_reports
WHERE id = $1;

-- name: UpsertCommunityReportVote :one
INSERT INTO community_report_votes (
    community_report_id,
    user_id,
    vote
) VALUES (
    $1,
    $2,
    $3
)
ON CONFLICT (community_report_id, user_id) DO UPDATE SET
    vote = EXCLUDED.vote,
    updated_at = NOW()
RETURNING *;

-- name: DeleteExpiredCommunityReports :execrows
DELETE FROM community_reports
WHERE created_at <= NOW() - INTERVAL '2 hours';
