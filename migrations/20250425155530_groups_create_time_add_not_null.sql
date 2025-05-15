-- +goose Up
-- +goose StatementBegin
ALTER TABLE groups ALTER COLUMN created_at SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE groups ALTER COLUMN created_at DROP NOT NULL;
-- +goose StatementEnd
