-- +goose Up
CREATE TABLE classes (
    id                BIGSERIAL   PRIMARY KEY,
    department_id     BIGINT      NOT NULL REFERENCES departments(id) ON DELETE RESTRICT,
    admission_session TEXT        NOT NULL CHECK (admission_session ~ '^\d{4}/\d{4}$'),
    current_level_id  BIGINT      NOT NULL REFERENCES levels(id) ON DELETE RESTRICT,
    name              TEXT        NOT NULL,
    invite_code       TEXT        NOT NULL UNIQUE,
    created_by        BIGINT      REFERENCES users(id) ON DELETE SET NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (department_id, admission_session)
);

-- +goose Down
DROP TABLE IF EXISTS classes;
