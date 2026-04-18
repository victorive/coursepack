-- +goose Up
CREATE TABLE level_members (
    id        BIGSERIAL   PRIMARY KEY,
    level_id  BIGINT      NOT NULL REFERENCES levels(id) ON DELETE CASCADE,
    user_id   BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (level_id, user_id)
);

-- +goose Down
DROP TABLE level_members;
