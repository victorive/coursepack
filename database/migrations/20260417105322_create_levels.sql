-- +goose Up
CREATE TABLE levels (
    id    BIGSERIAL PRIMARY KEY,
    value SMALLINT  NOT NULL UNIQUE CHECK (value IN (100, 200, 300, 400, 500, 600, 700, 800)),
    label TEXT      NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS levels;
