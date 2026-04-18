-- +goose Up
CREATE TABLE lectures (
    id            BIGSERIAL   PRIMARY KEY,
    course_id     BIGINT      NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    topic         TEXT        NOT NULL,
    description   TEXT,
    lecturer_name TEXT,
    week_number   SMALLINT,
    kind          TEXT        NOT NULL DEFAULT 'lecture' CHECK (kind IN ('lecture', 'test', 'exam', 'lab')),
    status        TEXT        NOT NULL DEFAULT 'published' CHECK (status IN ('published', 'archived')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE lectures;
