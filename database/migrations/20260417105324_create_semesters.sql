-- +goose Up
CREATE TABLE semesters (
    id         BIGSERIAL   PRIMARY KEY,
    class_id   BIGINT      NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    session    TEXT        NOT NULL CHECK (session ~ '^\d{4}/\d{4}$'),
    kind       TEXT        NOT NULL CHECK (kind IN ('first', 'second', 'summer')),
    is_current BOOLEAN     NOT NULL DEFAULT FALSE,
    starts_at  DATE,
    ends_at    DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (class_id, session, kind)
);

CREATE UNIQUE INDEX semesters_one_current_per_class
    ON semesters (class_id, is_current)
    WHERE is_current = TRUE;

-- +goose Down
DROP TABLE semesters;
