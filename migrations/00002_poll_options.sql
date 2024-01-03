-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS poll_options (
    id bigserial PRIMARY KEY,
    poll_id bigserial REFERENCES polls (id) ON DELETE CASCADE,
    value text NOT NULL,
    position int NOT NULL,  
    vote_count int NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS poll_options;
-- +goose StatementEnd
