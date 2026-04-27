-- +goose Up
CREATE TABLE class_courses (
    id          BIGSERIAL   PRIMARY KEY,
    class_id    BIGINT      NOT NULL REFERENCES classes(id)   ON DELETE CASCADE,
    semester_id BIGINT      NOT NULL REFERENCES semesters(id) ON DELETE CASCADE,
    course_id   BIGINT      NOT NULL REFERENCES courses(id)   ON DELETE RESTRICT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (class_id, semester_id, course_id)
);

-- +goose Down
DROP TABLE IF EXISTS class_courses;
