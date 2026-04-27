-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_points (
    user_id      BIGINT      PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    balance      INT         NOT NULL DEFAULT 0,
    total_earned INT         NOT NULL DEFAULT 0,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_points;
-- +goose StatementEnd
