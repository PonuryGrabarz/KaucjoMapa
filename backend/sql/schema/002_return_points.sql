-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS cube;
CREATE EXTENSION IF NOT EXISTS earthdistance;

CREATE TYPE point_type AS ENUM (
    'butelkomat',
    'sklep',
    'punkt_skupu',
    'inny'
);

CREATE TYPE point_status AS ENUM (
    'dostepny',
    'zapelniony',
    'nieczynny',
    'awaria',
    'brak_informacji',
    'do_weryfikacji'
);

CREATE TABLE return_points (
    id                  UUID         PRIMARY KEY,
    user_id    UUID         REFERENCES users (id) ON DELETE SET NULL,
    name                VARCHAR(255) NOT NULL,
    type                point_type   NOT NULL DEFAULT 'inny',
    address             TEXT,
    latitude            DOUBLE PRECISION NOT NULL,
    longitude           DOUBLE PRECISION NOT NULL,
    pet                 BOOLEAN NOT NULL DEFAULT FALSE,
    can                 BOOLEAN NOT NULL DEFAULT FALSE,
    glass               BOOLEAN NOT NULL DEFAULT FALSE,
    opening_hours       TEXT,
    description         TEXT,
    current_status      point_status NOT NULL DEFAULT 'brak_informacji',
    confirmations_count INT          NOT NULL DEFAULT 0,
    rejections_count    INT          NOT NULL DEFAULT 0,
    is_verified         BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_return_points_location ON return_points USING GIST (
    ll_to_earth(latitude, longitude)
);

CREATE UNIQUE INDEX uniq_return_points_location_rounded
ON return_points (
    ROUND(latitude::numeric, 5),
    ROUND(longitude::numeric, 5)
);

CREATE INDEX idx_return_points_status   ON return_points (current_status);
CREATE INDEX idx_return_points_type     ON return_points (type);
CREATE INDEX idx_return_points_user     ON return_points (user_id);

-- +goose Down
DROP TABLE IF EXISTS return_points;
