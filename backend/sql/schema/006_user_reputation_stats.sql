-- +goose Up
ALTER TABLE users
ADD COLUMN total_confirmations INT NOT NULL DEFAULT 0,
ADD COLUMN total_rejections INT NOT NULL DEFAULT 0,
ADD CONSTRAINT users_total_confirmations_nonnegative CHECK (total_confirmations >= 0),
ADD CONSTRAINT users_total_rejections_nonnegative CHECK (total_rejections >= 0);

DELETE FROM community_report_votes AS crv
USING community_reports AS cr
WHERE cr.id = crv.community_report_id
  AND cr.user_id = crv.user_id;

UPDATE community_reports AS cr
SET
    confirmations_count = (
        SELECT COUNT(*)
        FROM community_report_votes AS crv
        WHERE crv.community_report_id = cr.id
          AND crv.vote = 'confirm'
    ),
    rejections_count = (
        SELECT COUNT(*)
        FROM community_report_votes AS crv
        WHERE crv.community_report_id = cr.id
          AND crv.vote = 'reject'
    );

WITH user_vote_totals AS (
    SELECT
        u.id AS user_id,
        COUNT(crv.*) FILTER (WHERE crv.vote = 'confirm')::INT AS confirmations,
        COUNT(crv.*) FILTER (WHERE crv.vote = 'reject')::INT AS rejections
    FROM users AS u
    LEFT JOIN community_reports AS cr ON cr.user_id = u.id
    LEFT JOIN community_report_votes AS crv ON crv.community_report_id = cr.id
    GROUP BY u.id
)
UPDATE users AS u
SET
    total_confirmations = uvt.confirmations,
    total_rejections = uvt.rejections,
    reputation = CASE
        WHEN uvt.confirmations + uvt.rejections = 0 THEN 0
        ELSE ROUND(
            100.0 * uvt.confirmations / (uvt.confirmations + uvt.rejections)
        )::INT
    END
FROM user_vote_totals AS uvt
WHERE u.id = uvt.user_id;

ALTER TABLE users
ADD CONSTRAINT users_reputation_range CHECK (reputation BETWEEN 0 AND 100);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION apply_community_report_vote()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
DECLARE
    report_author_id UUID;
    confirmation_delta INT;
    rejection_delta INT;
    total_votes INT;
    rejection_votes INT;
BEGIN
    SELECT user_id
    INTO report_author_id
    FROM community_reports
    WHERE id = NEW.community_report_id;

    IF report_author_id = NEW.user_id THEN
        RAISE EXCEPTION 'users cannot vote on their own community reports'
            USING ERRCODE = 'check_violation';
    END IF;

    IF TG_OP = 'INSERT' THEN
        confirmation_delta := (NEW.vote = 'confirm')::INT;
        rejection_delta := (NEW.vote = 'reject')::INT;
    ELSE
        confirmation_delta := (NEW.vote = 'confirm')::INT - (OLD.vote = 'confirm')::INT;
        rejection_delta := (NEW.vote = 'reject')::INT - (OLD.vote = 'reject')::INT;
    END IF;

    UPDATE community_reports
    SET
        confirmations_count = confirmations_count + confirmation_delta,
        rejections_count = rejections_count + rejection_delta,
        updated_at = NOW()
    WHERE id = NEW.community_report_id
    RETURNING
        confirmations_count + rejections_count,
        rejections_count
    INTO total_votes, rejection_votes;

    UPDATE users
    SET
        total_confirmations = total_confirmations + confirmation_delta,
        total_rejections = total_rejections + rejection_delta,
        reputation = CASE
            WHEN total_confirmations + confirmation_delta
                + total_rejections + rejection_delta = 0 THEN 0
            ELSE ROUND(
                100.0 * (total_confirmations + confirmation_delta)
                / (
                    total_confirmations + confirmation_delta
                    + total_rejections + rejection_delta
                )
            )::INT
        END,
        updated_at = NOW()
    WHERE id = report_author_id;

    IF total_votes >= 10 AND rejection_votes * 4 >= total_votes * 3 THEN
        DELETE FROM community_reports
        WHERE id = NEW.community_report_id;
    END IF;

    RETURN NEW;
END;
$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION apply_community_report_vote()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
DECLARE
    total_votes INT;
    rejection_votes INT;
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE community_reports
        SET
            confirmations_count = confirmations_count + (NEW.vote = 'confirm')::INT,
            rejections_count = rejections_count + (NEW.vote = 'reject')::INT,
            updated_at = NOW()
        WHERE id = NEW.community_report_id
        RETURNING
            confirmations_count + rejections_count,
            rejections_count
        INTO total_votes, rejection_votes;
    ELSE
        UPDATE community_reports
        SET
            confirmations_count = confirmations_count
                - (OLD.vote = 'confirm')::INT
                + (NEW.vote = 'confirm')::INT,
            rejections_count = rejections_count
                - (OLD.vote = 'reject')::INT
                + (NEW.vote = 'reject')::INT,
            updated_at = NOW()
        WHERE id = NEW.community_report_id
        RETURNING
            confirmations_count + rejections_count,
            rejections_count
        INTO total_votes, rejection_votes;
    END IF;

    IF total_votes >= 10 AND rejection_votes * 4 >= total_votes * 3 THEN
        DELETE FROM community_reports
        WHERE id = NEW.community_report_id;
    END IF;

    RETURN NEW;
END;
$$;
-- +goose StatementEnd

ALTER TABLE users
DROP CONSTRAINT IF EXISTS users_reputation_range,
DROP CONSTRAINT IF EXISTS users_total_rejections_nonnegative,
DROP CONSTRAINT IF EXISTS users_total_confirmations_nonnegative,
DROP COLUMN IF EXISTS total_rejections,
DROP COLUMN IF EXISTS total_confirmations;
