-- +goose Up
CREATE TYPE data_error_category AS ENUM (
    'adres_lokalizacja',
    'opis_punktu',
    'formy_rozliczenia',
    'godziny_otwarcia',
    'przyjmowane_opakowania',
    'inne_dane'
);

CREATE TYPE data_error_report_status AS ENUM (
    'oczekujace',
    'rozpatrzone',
    'odrzucone'
);

CREATE TABLE data_error_reports (
    id                   UUID                      PRIMARY KEY,
    user_id              UUID                      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    return_point_id      UUID                      NOT NULL REFERENCES return_points (id) ON DELETE CASCADE,
    categories           data_error_category[]     NOT NULL,
    description          VARCHAR(500)              NOT NULL,
    suggested_correction TEXT,
    status               data_error_report_status  NOT NULL DEFAULT 'oczekujace',
    created_at           TIMESTAMPTZ               NOT NULL DEFAULT NOW(),
    updated_at           TIMESTAMPTZ               NOT NULL DEFAULT NOW(),

    CONSTRAINT data_error_reports_categories_not_empty
        CHECK (cardinality(categories) > 0),
    CONSTRAINT data_error_reports_description_not_blank
        CHECK (char_length(trim(description)) > 0)
);

CREATE INDEX idx_data_error_reports_status_created_at
ON data_error_reports (status, created_at DESC);

CREATE INDEX idx_data_error_reports_return_point_created_at
ON data_error_reports (return_point_id, created_at DESC);

CREATE INDEX idx_data_error_reports_user_id
ON data_error_reports (user_id);

-- +goose Down
DROP TABLE IF EXISTS data_error_reports;
DROP TYPE IF EXISTS data_error_report_status;
DROP TYPE IF EXISTS data_error_category;
