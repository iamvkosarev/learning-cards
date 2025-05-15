-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
id bigint PRIMARY KEY NOT NULL);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd