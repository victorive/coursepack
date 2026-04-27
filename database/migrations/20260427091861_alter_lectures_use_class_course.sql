-- +goose Up
ALTER TABLE lectures DROP CONSTRAINT IF EXISTS lectures_course_id_fkey;
ALTER TABLE lectures ADD COLUMN class_course_id BIGINT REFERENCES class_courses(id) ON DELETE CASCADE;
ALTER TABLE lectures ALTER COLUMN class_course_id SET NOT NULL;
ALTER TABLE lectures DROP COLUMN course_id;

-- +goose Down
ALTER TABLE lectures ADD COLUMN course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE;
ALTER TABLE lectures DROP COLUMN IF EXISTS class_course_id;
