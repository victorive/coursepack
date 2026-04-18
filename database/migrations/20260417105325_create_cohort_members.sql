-- +goose Up
CREATE TABLE cohort_members (
    id        BIGSERIAL   PRIMARY KEY,
    cohort_id BIGINT      NOT NULL REFERENCES cohorts(id) ON DELETE CASCADE,
    user_id   BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (cohort_id, user_id)
);

-- +goose Down
DROP TABLE cohort_members;
