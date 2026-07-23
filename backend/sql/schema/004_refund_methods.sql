-- +goose Up
CREATE TYPE refund_method AS ENUM (
    'gotowka',
    'bon',
    'karta'
);

ALTER TABLE return_points
ADD COLUMN refund_methods refund_method[]
NOT NULL DEFAULT '{}'::refund_method[];

-- +goose Down
ALTER TABLE return_points
DROP COLUMN IF EXISTS refund_methods;

DROP TYPE IF EXISTS refund_method;
Po dodaniu migracji uruchom:
