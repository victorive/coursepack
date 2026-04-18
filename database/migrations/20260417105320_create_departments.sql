-- +goose Up
CREATE TABLE departments (
    id         BIGSERIAL   PRIMARY KEY,
    name       TEXT        NOT NULL,
    code       TEXT        NOT NULL,
    faculty    TEXT        NOT NULL,
    university TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (code, university)
);

-- +goose Down
DROP TABLE departments;
