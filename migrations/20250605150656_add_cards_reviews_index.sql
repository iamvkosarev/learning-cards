-- +goose Up
-- +goose StatementBegin

CREATE INDEX IF NOT EXISTS idx_card_reviews_user_group_card_time
    ON card_reviews(user_id, group_id, card_id, time);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_card_reviews_user_group_card_time;
-- +goose StatementEnd
