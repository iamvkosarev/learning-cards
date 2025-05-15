-- +goose Up
-- +goose StatementBegin
ALTER TABLE cards RENAME CONSTRAINT cards_group_id_fkey TO  fk_cards_group_id_groups;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE cards RENAME CONSTRAINT fk_cards_group_id_groups TO  cards_group_id_fkey;
-- +goose StatementEnd
