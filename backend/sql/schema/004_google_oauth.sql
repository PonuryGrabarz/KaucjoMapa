-- +goose Up
ALTER TABLE users
ADD COLUMN google_subject TEXT;

CREATE UNIQUE INDEX users_google_subject_uidx
ON users (google_subject)
WHERE google_subject IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS users_google_subject_uidx;

ALTER TABLE users
DROP COLUMN google_subject;
