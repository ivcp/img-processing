-- +goose Up
-- +goose StatementBegin
ALTER TABLE polls ADD COLUMN results_visibility text NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE polls DROP COLUMN results_visibility;
-- +goose StatementEnd
