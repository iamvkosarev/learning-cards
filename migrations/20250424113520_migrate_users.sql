-- +goose Up
-- +goose StatementBegin
INSERT INTO users (id) SELECT user_id FROM card_groups GROUP BY user_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE id IN (SELECT user_id FROM card_groups);
-- +goose StatementEnd
