-- +goose Up
CREATE TABLE user_roles (
    id         BIGSERIAL   PRIMARY KEY,
    user_id    BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id    BIGINT      NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    scope_type TEXT        NOT NULL DEFAULT 'global' CHECK (scope_type IN ('global', 'class')),
    scope_id   BIGINT
);

CREATE UNIQUE INDEX user_roles_unique_global
    ON user_roles (user_id, role_id)
    WHERE scope_type = 'global';

CREATE UNIQUE INDEX user_roles_unique_class
    ON user_roles (user_id, role_id, scope_id)
    WHERE scope_type = 'class';

-- +goose Down
DROP TABLE user_roles;
