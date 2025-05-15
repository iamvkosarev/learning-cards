-- +goose Up
-- +goose StatementBegin
ALTER TABLE groups ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT now();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE groups DROP COLUMN IF EXISTS updated_at;
-- +goose StatementEnd
