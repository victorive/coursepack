-- +goose Up
-- +goose StatementBegin
ALTER TABLE departments DROP COLUMN IF EXISTS university;
ALTER TABLE departments ADD COLUMN university_id BIGINT NOT NULL;
ALTER TABLE departments ADD CONSTRAINT fk_departments_university_id FOREIGN KEY (university_id) REFERENCES universities(id) ON DELETE CASCADE;
ALTER TABLE departments DROP CONSTRAINT IF EXISTS "unique_code";
ALTER TABLE departments DROP CONSTRAINT IF EXISTS "departments_code_key";
ALTER TABLE departments ADD CONSTRAINT unique_code_university_id UNIQUE (code, university_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE departments CASCADE;
ALTER TABLE departments DROP CONSTRAINT IF EXISTS fk_departments_university_id;
ALTER TABLE departments DROP CONSTRAINT IF EXISTS unique_code_university_id;
ALTER TABLE departments DROP COLUMN IF EXISTS university_id;
ALTER TABLE departments ADD COLUMN university VARCHAR(255);
ALTER TABLE departments ADD CONSTRAINT "unique_code" UNIQUE (code);
-- +goose StatementEnd
