-- +goose Up
-- +goose StatementBegin
ALTER TABLE groups ALTER COLUMN visibility SET DATA TYPE numeric(1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE groups ALTER COLUMN visibility SET DATA TYPE smallint;
-- +goose StatementEnd
