-- +goose Up
CREATE TABLE class_members (
    id        BIGSERIAL   PRIMARY KEY,
    class_id  BIGINT      NOT NULL REFERENCES classes(id) ON DELETE CASCADE,
    user_id   BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (class_id, user_id)
);

-- +goose Down
DROP TABLE class_members;
