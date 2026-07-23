-- +goose Up
CREATE TYPE community_report_vote_value AS ENUM (
    'confirm',
    'reject'
);

CREATE TABLE community_reports (
    id                  UUID        PRIMARY KEY,
    user_id             UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    return_point_id     UUID        NOT NULL REFERENCES return_points (id) ON DELETE CASCADE,
    status              point_status NOT NULL DEFAULT 'brak_informacji',
    description         TEXT        NOT NULL,
    confirmations_count INT         NOT NULL DEFAULT 0 CHECK (confirmations_count >= 0),
    rejections_count    INT         NOT NULL DEFAULT 0 CHECK (rejections_count >= 0),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_community_reports_user_id
ON community_reports (user_id);

CREATE INDEX idx_community_reports_created_at
ON community_reports (created_at);

CREATE INDEX idx_community_reports_point_created_at
ON community_reports (return_point_id, created_at DESC);

CREATE TABLE community_report_votes (
    community_report_id UUID                  NOT NULL REFERENCES community_reports (id) ON DELETE CASCADE,
    user_id              UUID                  NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    vote                 community_report_vote_value NOT NULL,
    created_at           TIMESTAMPTZ           NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ           NOT NULL DEFAULT NOW(),
    PRIMARY KEY (community_report_id, user_id)
);

CREATE INDEX idx_community_report_votes_user_id
ON community_report_votes (user_id);

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

CREATE TRIGGER community_report_vote_inserted
AFTER INSERT ON community_report_votes
FOR EACH ROW
EXECUTE FUNCTION apply_community_report_vote();

CREATE TRIGGER community_report_vote_changed
AFTER UPDATE OF vote ON community_report_votes
FOR EACH ROW
WHEN (OLD.vote IS DISTINCT FROM NEW.vote)
EXECUTE FUNCTION apply_community_report_vote();

-- +goose Down
DROP TRIGGER IF EXISTS community_report_vote_changed ON community_report_votes;
DROP TRIGGER IF EXISTS community_report_vote_inserted ON community_report_votes;
DROP FUNCTION IF EXISTS apply_community_report_vote();
DROP TABLE IF EXISTS community_report_votes;
DROP TABLE IF EXISTS community_reports;
DROP TYPE IF EXISTS community_report_vote_value;
