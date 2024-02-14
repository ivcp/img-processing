-- +goose Up
-- +goose StatementBegin
ALTER TABLE polls ADD COLUMN is_private BOOLEAN NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE polls DROP COLUMN is_private;
-- +goose StatementEnd
