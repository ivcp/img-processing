-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS polls (
    id bigserial PRIMARY KEY,
    question text NOT NULL,
    description text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    expires_at timestamp(0) with time zone NOT NULL,
    version int NOT NULL DEFAULT 1
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS polls;
-- +goose StatementEnd
