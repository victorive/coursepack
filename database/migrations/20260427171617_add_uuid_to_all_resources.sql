-- +goose Up
ALTER TABLE universities ADD COLUMN uuid VARCHAR(36) NOT NULL DEFAULT '';
UPDATE universities SET uuid = gen_random_uuid()::text WHERE uuid = '';
ALTER TABLE universities ALTER COLUMN uuid DROP DEFAULT;
ALTER TABLE universities ADD CONSTRAINT universities_uuid_unique UNIQUE (uuid);

ALTER TABLE departments ADD COLUMN uuid VARCHAR(36) NOT NULL DEFAULT '';
UPDATE departments SET uuid = gen_random_uuid()::text WHERE uuid = '';
ALTER TABLE departments ALTER COLUMN uuid DROP DEFAULT;
ALTER TABLE departments ADD CONSTRAINT departments_uuid_unique UNIQUE (uuid);

ALTER TABLE classes ADD COLUMN uuid VARCHAR(36) NOT NULL DEFAULT '';
UPDATE classes SET uuid = gen_random_uuid()::text WHERE uuid = '';
ALTER TABLE classes ALTER COLUMN uuid DROP DEFAULT;
ALTER TABLE classes ADD CONSTRAINT classes_uuid_unique UNIQUE (uuid);

ALTER TABLE courses ADD COLUMN uuid VARCHAR(36) NOT NULL DEFAULT '';
UPDATE courses SET uuid = gen_random_uuid()::text WHERE uuid = '';
ALTER TABLE courses ALTER COLUMN uuid DROP DEFAULT;
ALTER TABLE courses ADD CONSTRAINT courses_uuid_unique UNIQUE (uuid);

ALTER TABLE class_courses ADD COLUMN uuid VARCHAR(36) NOT NULL DEFAULT '';
UPDATE class_courses SET uuid = gen_random_uuid()::text WHERE uuid = '';
ALTER TABLE class_courses ALTER COLUMN uuid DROP DEFAULT;
ALTER TABLE class_courses ADD CONSTRAINT class_courses_uuid_unique UNIQUE (uuid);

ALTER TABLE lectures ADD COLUMN uuid VARCHAR(36) NOT NULL DEFAULT '';
UPDATE lectures SET uuid = gen_random_uuid()::text WHERE uuid = '';
ALTER TABLE lectures ALTER COLUMN uuid DROP DEFAULT;
ALTER TABLE lectures ADD CONSTRAINT lectures_uuid_unique UNIQUE (uuid);

ALTER TABLE files ADD COLUMN uuid VARCHAR(36) NOT NULL DEFAULT '';
UPDATE files SET uuid = gen_random_uuid()::text WHERE uuid = '';
ALTER TABLE files ALTER COLUMN uuid DROP DEFAULT;
ALTER TABLE files ADD CONSTRAINT files_uuid_unique UNIQUE (uuid);

-- +goose Down
ALTER TABLE universities DROP CONSTRAINT IF EXISTS universities_uuid_unique;
ALTER TABLE universities DROP COLUMN IF EXISTS uuid;

ALTER TABLE departments DROP CONSTRAINT IF EXISTS departments_uuid_unique;
ALTER TABLE departments DROP COLUMN IF EXISTS uuid;

ALTER TABLE classes DROP CONSTRAINT IF EXISTS classes_uuid_unique;
ALTER TABLE classes DROP COLUMN IF EXISTS uuid;

ALTER TABLE courses DROP CONSTRAINT IF EXISTS courses_uuid_unique;
ALTER TABLE courses DROP COLUMN IF EXISTS uuid;

ALTER TABLE class_courses DROP CONSTRAINT IF EXISTS class_courses_uuid_unique;
ALTER TABLE class_courses DROP COLUMN IF EXISTS uuid;

ALTER TABLE lectures DROP CONSTRAINT IF EXISTS lectures_uuid_unique;
ALTER TABLE lectures DROP COLUMN IF EXISTS uuid;

ALTER TABLE files DROP CONSTRAINT IF EXISTS files_uuid_unique;
ALTER TABLE files DROP COLUMN IF EXISTS uuid;
