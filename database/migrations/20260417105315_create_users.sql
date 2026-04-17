-- +goose Up
CREATE TABLE users (
    id             BIGSERIAL   PRIMARY KEY,
    uuid           UUID        NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    email          TEXT        NOT NULL UNIQUE,
    password       TEXT,
    full_name      TEXT        NOT NULL,
    phone          TEXT,
    avatar_url     TEXT,
    provider       TEXT,
    email_verified BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE users;
