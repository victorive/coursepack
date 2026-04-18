-- +goose Up
CREATE TABLE levels (
    id            BIGSERIAL   PRIMARY KEY,
    department_id BIGINT      NOT NULL REFERENCES departments(id) ON DELETE CASCADE,
    year          SMALLINT    NOT NULL CHECK (year IN (100, 200, 300, 400, 500, 600, 700, 800)),
    invite_code   TEXT        NOT NULL UNIQUE,
    created_by    BIGINT      REFERENCES users(id) ON DELETE SET NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (department_id, year)
);

-- +goose Down
DROP TABLE levels;
