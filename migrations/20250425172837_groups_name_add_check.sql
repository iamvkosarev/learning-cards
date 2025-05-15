-- +goose Up
-- +goose StatementBegin
ALTER TABLE groups ADD CONSTRAINT name_min_size CHECK ( length(name) > 0 );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE groups DROP CONSTRAINT IF EXISTS name_min_size;
-- +goose StatementEnd
