-- +goose Up
-- +goose StatementBegin
DELETE FROM point_transactions a
USING point_transactions b
WHERE a.id > b.id
  AND a.source_id IS NOT NULL
  AND b.source_id IS NOT NULL
  AND a.user_id = b.user_id
  AND a.source_type = b.source_type
  AND a.source_id = b.source_id;

CREATE UNIQUE INDEX uq_point_transactions_source
    ON point_transactions (user_id, source_type, source_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS uq_point_transactions_source;
-- +goose StatementEnd
