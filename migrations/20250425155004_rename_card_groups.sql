-- +goose Up
-- +goose StatementBegin
ALTER TABLE card_groups RENAME TO groups;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE groups RENAME TO card_groups;
-- +goose StatementEnd
