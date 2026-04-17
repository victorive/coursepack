-- +goose Up
CREATE TABLE user_roles (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id    BIGINT      NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    scope_type TEXT        NOT NULL DEFAULT 'global' CHECK (scope_type IN ('global', 'cohort')),
    scope_id   BIGINT,
    granted_by BIGINT      REFERENCES users(id) ON DELETE SET NULL,
    granted_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX user_roles_unique_global
    ON user_roles (user_id, role_id)
    WHERE scope_type = 'global';

CREATE UNIQUE INDEX user_roles_unique_cohort
    ON user_roles (user_id, role_id, scope_id)
    WHERE scope_type = 'cohort';

-- +goose Down
DROP TABLE user_roles;
