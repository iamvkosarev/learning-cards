-- +goose Up
-- +goose StatementBegin
ALTER TABLE groups DROP CONSTRAINT IF EXISTS fk_card_groups_user_id_users;

ALTER TABLE groups ADD CONSTRAINT fk_groups_user_id_users
FOREIGN KEY (user_id) REFERENCES users(id)
ON UPDATE CASCADE
ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE groups DROP CONSTRAINT IF EXISTS fk_groups_user_id_users;

ALTER TABLE groups ADD CONSTRAINT fk_card_groups_user_id_users
FOREIGN KEY (user_id) REFERENCES users(id)
ON DELETE CASCADE;
-- +goose StatementEnd
