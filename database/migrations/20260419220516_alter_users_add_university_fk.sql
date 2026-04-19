-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN university_id BIGINT;
ALTER TABLE users ADD CONSTRAINT fk_users_university_id FOREIGN KEY (university_id) REFERENCES universities(id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_university_id;
ALTER TABLE users DROP COLUMN IF EXISTS university_id;
-- +goose StatementEnd
