-- +goose Up
-- +goose StatementBegin
ALTER TABLE card_groups
	ADD CONSTRAINT fk_card_groups_user_id_users
	FOREIGN KEY (user_id) REFERENCES users(id)
	ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE card_groups
	DROP CONSTRAINT IF EXISTS fk_card_groups_user_id_users;
-- +goose StatementEnd
