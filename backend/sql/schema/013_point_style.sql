-- +goose Up
ALTER TABLE return_points
ADD COLUMN short_name TEXT,
ADD COLUMN bg_color   TEXT,
ADD COLUMN text_color TEXT;

-- +goose Down
ALTER TABLE return_points
DROP COLUMN short_name,
DROP COLUMN bg_color,
DROP COLUMN text_color;
