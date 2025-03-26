-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS card_groups (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT NOT NULL,
	name TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT now()
	);

CREATE TABLE IF NOT EXISTS cards (
	id BIGSERIAL PRIMARY KEY,
	group_id BIGINT NOT NULL REFERENCES card_groups (id) ON DELETE CASCADE,
	front_text TEXT NOT NULL,
	back_text TEXT NOT NULL,
	created_at TIMESTAMPTZ DEFAULT now()
	);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cards;

DROP TABLE IF EXISTS card_groups;

-- +goose StatementEnd