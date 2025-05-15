-- +goose Up
-- +goose StatementBegin
ALTER TABLE cards ALTER COLUMN front_text SET DATA TYPE varchar(250);
ALTER TABLE cards ALTER COLUMN back_text SET DATA TYPE varchar(250);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE cards ALTER COLUMN front_text SET DATA TYPE text;
ALTER TABLE cards ALTER COLUMN back_text SET DATA TYPE text;
-- +goose StatementEnd
