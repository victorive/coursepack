-- +goose Up
CREATE TABLE departments (
    id         BIGSERIAL   PRIMARY KEY,
    name       TEXT        NOT NULL,
    code       TEXT        NOT NULL,
    faculty    TEXT,
    university TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (code)
);

-- +goose Down
DROP TABLE departments;
