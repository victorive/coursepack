-- +goose Up
-- +goose StatementBegin
CREATE TABLE point_transactions (
    id            BIGSERIAL   PRIMARY KEY,
    user_id       BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    delta         INT         NOT NULL,
    balance_after INT         NOT NULL,
    source_type   TEXT        NOT NULL CHECK (source_type IN ('file_approval', 'manual', 'bonus')),
    source_id     BIGINT,
    note          TEXT,
    created_by    BIGINT      REFERENCES users(id) ON DELETE SET NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_point_transactions_user ON point_transactions (user_id);
CREATE INDEX idx_point_transactions_source ON point_transactions (source_type, source_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS point_transactions;
-- +goose StatementEnd
