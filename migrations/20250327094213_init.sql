-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS card_groups (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    visibility SMALLINT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_card_groups_user_id ON card_groups (user_id);

CREATE TABLE IF NOT EXISTS cards (
    id BIGSERIAL PRIMARY KEY,
    group_id BIGINT NOT NULL REFERENCES card_groups (id) ON DELETE CASCADE,
    front_text TEXT NOT NULL,
    back_text TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_cards_group_id ON cards (group_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_cards_group_id;
DROP INDEX IF EXISTS idx_card_groups_user_id;

DROP TABLE IF EXISTS cards;
DROP TABLE IF EXISTS card_groups;

-- +goose StatementEnd