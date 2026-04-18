-- +goose Up
CREATE TABLE files (
    id          BIGSERIAL   PRIMARY KEY,
    uploaded_by BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id   BIGINT      NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    lecture_id  BIGINT      REFERENCES lectures(id) ON DELETE SET NULL,
    category    TEXT        NOT NULL CHECK (category IN ('lecture_note', 'past_question', 'material', 'assignment')),
    title       TEXT        NOT NULL,
    storage_key TEXT        NOT NULL UNIQUE,
    size_bytes  BIGINT      NOT NULL,
    mime_type   TEXT        NOT NULL,
    year_label  TEXT,
    status      TEXT        NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'archived')),
    reviewed_by BIGINT      REFERENCES users(id) ON DELETE SET NULL,
    review_note TEXT,
    reviewed_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE files;
