-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX uq_point_transactions_source
    ON point_transactions (user_id, source_type, source_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS uq_point_transactions_source;
-- +goose StatementEnd
