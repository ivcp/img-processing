-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ips (   
    id bigserial PRIMARY KEY,
    ip inet NOT NULL,
    poll_id bigserial REFERENCES polls (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ips;
-- +goose StatementEnd
