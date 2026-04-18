-- +goose Up
CREATE TABLE academic_sessions (
    id         BIGSERIAL   PRIMARY KEY,
    label      TEXT        NOT NULL UNIQUE,
    is_current BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX academic_sessions_one_current
    ON academic_sessions (is_current)
    WHERE is_current = TRUE;

-- +goose Down
DROP TABLE academic_sessions;
