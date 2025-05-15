-- +goose Up
-- +goose StatementBegin
ALTER TABLE cards ALTER COLUMN created_at SET NOT NULL;
ALTER TABLE cards ALTER COLUMN created_at SET DATA TYPE timestamp;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE cards ALTER COLUMN created_at DROP NOT NULL;
ALTER TABLE cards ALTER COLUMN created_at SET DATA TYPE timestamptz;
-- +goose StatementEnd
