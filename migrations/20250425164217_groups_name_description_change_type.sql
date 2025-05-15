-- +goose Up
-- +goose StatementBegin
ALTER TABLE groups ALTER COLUMN name SET DATA TYPE varchar(100);
ALTER TABLE groups ALTER COLUMN description SET DATA TYPE varchar(300);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE groups ALTER COLUMN name SET DATA TYPE text;
ALTER TABLE groups ALTER COLUMN description SET DATA TYPE text;
-- +goose StatementEnd
