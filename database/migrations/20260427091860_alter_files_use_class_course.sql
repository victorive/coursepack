-- +goose Up
ALTER TABLE files DROP CONSTRAINT IF EXISTS files_course_id_fkey;
ALTER TABLE files ADD COLUMN class_course_id BIGINT REFERENCES class_courses(id) ON DELETE CASCADE;
ALTER TABLE files ALTER COLUMN class_course_id SET NOT NULL;
ALTER TABLE files DROP COLUMN course_id;

-- +goose Down
ALTER TABLE files ADD COLUMN course_id BIGINT NOT NULL REFERENCES courses(id) ON DELETE CASCADE;
ALTER TABLE files DROP COLUMN IF EXISTS class_course_id;
