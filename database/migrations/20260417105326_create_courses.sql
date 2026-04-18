-- +goose Up
CREATE TABLE courses (
    id          BIGSERIAL   PRIMARY KEY,
    cohort_id   BIGINT      NOT NULL REFERENCES cohorts(id) ON DELETE CASCADE,
    semester_id BIGINT      NOT NULL REFERENCES semesters(id) ON DELETE CASCADE,
    code        TEXT        NOT NULL,
    title       TEXT        NOT NULL,
    units       SMALLINT    NOT NULL DEFAULT 2,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (cohort_id, semester_id, code)
);

-- +goose Down
DROP TABLE courses;
