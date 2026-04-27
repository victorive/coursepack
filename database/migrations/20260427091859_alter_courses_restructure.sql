-- +goose Up
ALTER TABLE courses DROP CONSTRAINT IF EXISTS courses_class_id_semester_id_code_key;
ALTER TABLE courses DROP CONSTRAINT IF EXISTS courses_class_id_fkey;
ALTER TABLE courses DROP CONSTRAINT IF EXISTS courses_semester_id_fkey;

ALTER TABLE courses
    ADD COLUMN university_id BIGINT REFERENCES universities(id) ON DELETE RESTRICT,
    ADD COLUMN department_id BIGINT REFERENCES departments(id)  ON DELETE RESTRICT,
    ADD COLUMN credit_units  SMALLINT NOT NULL DEFAULT 3;

ALTER TABLE courses
    DROP COLUMN class_id,
    DROP COLUMN semester_id,
    DROP COLUMN units;

ALTER TABLE courses ALTER COLUMN university_id SET NOT NULL;
ALTER TABLE courses ALTER COLUMN department_id SET NOT NULL;

ALTER TABLE courses ADD CONSTRAINT courses_university_department_code_key
    UNIQUE (university_id, department_id, code);

-- +goose Down
ALTER TABLE courses DROP CONSTRAINT IF EXISTS courses_university_department_code_key;
ALTER TABLE courses ALTER COLUMN university_id DROP NOT NULL;
ALTER TABLE courses ALTER COLUMN department_id DROP NOT NULL;
ALTER TABLE courses
    DROP COLUMN IF EXISTS university_id,
    DROP COLUMN IF EXISTS department_id,
    DROP COLUMN IF EXISTS credit_units;
ALTER TABLE courses
    ADD COLUMN class_id    BIGINT NOT NULL REFERENCES classes(id)   ON DELETE CASCADE,
    ADD COLUMN semester_id BIGINT NOT NULL REFERENCES semesters(id) ON DELETE CASCADE,
    ADD COLUMN units       SMALLINT NOT NULL DEFAULT 2;
ALTER TABLE courses ADD CONSTRAINT courses_class_id_semester_id_code_key
    UNIQUE (class_id, semester_id, code);
