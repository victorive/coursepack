-- +goose Up
CREATE TABLE join_requests (
    id           BIGSERIAL   PRIMARY KEY,
    level_id     BIGINT      NOT NULL REFERENCES levels(id) ON DELETE CASCADE,
    user_id      BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status       TEXT        NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    reviewed_by  BIGINT      REFERENCES users(id) ON DELETE SET NULL,
    requested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    resolved_at  TIMESTAMPTZ,
    UNIQUE (level_id, user_id)
);

-- +goose Down
DROP TABLE join_requests;
