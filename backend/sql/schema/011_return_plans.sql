-- +goose Up
CREATE TYPE return_plan_status AS ENUM (
    'draft',
    'saved',
    'partially_confirmed',
    'confirmed'
);

CREATE TYPE return_plan_point_status AS ENUM (
    'pending',
    'confirmed'
);

CREATE TABLE return_plans (
    id                 UUID               PRIMARY KEY,
    user_id            UUID               NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    status             return_plan_status NOT NULL DEFAULT 'draft',
    target_pet         INT                NOT NULL DEFAULT 0 CHECK (target_pet >= 0),
    target_can         INT                NOT NULL DEFAULT 0 CHECK (target_can >= 0),
    target_glass       INT                NOT NULL DEFAULT 0 CHECK (target_glass >= 0),
    rate_pet_grosze    INT                NOT NULL DEFAULT 50 CHECK (rate_pet_grosze >= 0),
    rate_can_grosze    INT                NOT NULL DEFAULT 50 CHECK (rate_can_grosze >= 0),
    rate_glass_grosze  INT                NOT NULL DEFAULT 100 CHECK (rate_glass_grosze >= 0),
    confirmed_at       TIMESTAMPTZ,
    created_at         TIMESTAMPTZ        NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ        NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_return_plans_user_created_at
ON return_plans (user_id, created_at DESC);

CREATE UNIQUE INDEX uniq_return_plans_one_draft_per_user
ON return_plans (user_id)
WHERE status = 'draft';

CREATE TABLE return_plan_points (
    id              UUID                      PRIMARY KEY,
    return_plan_id  UUID                      NOT NULL REFERENCES return_plans (id) ON DELETE CASCADE,
    return_point_id UUID                      NOT NULL REFERENCES return_points (id) ON DELETE RESTRICT,
    point_name      VARCHAR(255)              NOT NULL,
    point_address   TEXT,
    pet_qty         INT                       NOT NULL DEFAULT 0 CHECK (pet_qty >= 0),
    can_qty         INT                       NOT NULL DEFAULT 0 CHECK (can_qty >= 0),
    glass_qty       INT                       NOT NULL DEFAULT 0 CHECK (glass_qty >= 0),
    status          return_plan_point_status  NOT NULL DEFAULT 'pending',
    confirmed_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ               NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ               NOT NULL DEFAULT NOW(),

    CONSTRAINT uniq_return_plan_points_plan_point
        UNIQUE (return_plan_id, return_point_id),
    CONSTRAINT chk_return_plan_points_confirmed_at CHECK (
        (status = 'pending' AND confirmed_at IS NULL)
        OR (status = 'confirmed' AND confirmed_at IS NOT NULL)
    )
);

CREATE INDEX idx_return_plan_points_plan_id
ON return_plan_points (return_plan_id);

-- +goose Down
DROP TABLE IF EXISTS return_plan_points;
DROP TABLE IF EXISTS return_plans;
DROP TYPE IF EXISTS return_plan_point_status;
DROP TYPE IF EXISTS return_plan_status;
