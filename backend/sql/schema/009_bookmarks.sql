-- +goose Up
CREATE TABLE bookmarks (
    id              UUID        PRIMARY KEY,
    user_id         UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    return_point_id UUID        NOT NULL REFERENCES return_points (id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_bookmarks_user_point UNIQUE (user_id, return_point_id)
);

CREATE INDEX idx_bookmarks_user_created_at
ON bookmarks (user_id, created_at DESC);

CREATE INDEX idx_bookmarks_return_point_id
ON bookmarks (return_point_id);

-- +goose Down
DROP TABLE IF EXISTS bookmarks;
