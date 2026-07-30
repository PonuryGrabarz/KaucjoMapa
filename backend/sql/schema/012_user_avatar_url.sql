-- +goose Up
ALTER TABLE users
ADD COLUMN avatar_path TEXT;

-- +goose Down
ALTER TABLE users
DROP COLUMN avatar_path;
