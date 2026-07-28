-- +goose Up
CREATE TABLE password_reset_challenges (
    id                     UUID PRIMARY KEY,
    user_id                UUID,
    code_hash              TEXT        NOT NULL,
    attempts               INT         NOT NULL DEFAULT 0,
    expires_at             TIMESTAMPTZ NOT NULL,
    resend_available_at    TIMESTAMPTZ NOT NULL,
    reset_token_hash       TEXT,
    reset_token_expires_at TIMESTAMPTZ,
    consumed_at            TIMESTAMPTZ,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_password_reset_challenges_user_id ON password_reset_challenges (user_id);

-- +goose Down
DROP TABLE password_reset_challenges;
