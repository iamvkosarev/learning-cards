-- +goose Up
-- +goose StatementBegin
ALTER TABLE cards ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP NOT NULL DEFAULT now();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE cards DROP COLUMN IF EXISTS updated_at;
-- +goose StatementEnd
