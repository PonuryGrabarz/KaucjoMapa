-- +goose Up
CREATE TABLE user_notifications (
    id          UUID        PRIMARY KEY,
    user_id     UUID        REFERENCES users (id) ON DELETE CASCADE,
    name        VARCHAR(255) NOT NULL,
    description TEXT        NOT NULL,
    seen        BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_user_notifications_user_id
ON user_notifications (user_id);

CREATE INDEX idx_user_notifications_user_created_at
ON user_notifications (user_id, created_at DESC);

-- +goose Down
DROP TABLE IF EXISTS user_notifications;
