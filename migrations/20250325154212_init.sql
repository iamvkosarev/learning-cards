-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS card_groups (
                                           id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name TEXT NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now()
    );
CREATE TABLE IF NOT EXISTS cards (
                                     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    group_id UUID NOT NULL REFERENCES card_groups(id) ON DELETE CASCADE,
    front_text TEXT NOT NULL,
    back_text TEXT NOT NULL
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cards;
DROP TABLE IF EXISTS card_groups;
-- +goose StatementEnd
