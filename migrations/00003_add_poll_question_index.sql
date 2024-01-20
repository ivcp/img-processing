-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS polls_question_idx ON polls USING GIN (to_tsvector('simple', question));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS  polls_question_idx;
-- +goose StatementEnd
