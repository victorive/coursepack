-- +goose Up
CREATE TABLE semesters (
    id         BIGSERIAL   PRIMARY KEY,
    session_id BIGINT      NOT NULL REFERENCES academic_sessions(id) ON DELETE CASCADE,
    level_id   BIGINT      REFERENCES levels(id) ON DELETE CASCADE,
    kind       TEXT        NOT NULL CHECK (kind IN ('first', 'second', 'summer')),
    is_current BOOLEAN     NOT NULL DEFAULT FALSE,
    starts_at  DATE,
    ends_at    DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (level_id, session_id, kind)
);

CREATE UNIQUE INDEX semesters_one_current_per_level
    ON semesters (level_id, is_current)
    WHERE is_current = TRUE;

-- +goose Down
DROP TABLE semesters;
