-- +goose Up
-- +goose StatementBegin

ALTER TABLE cards
	RENAME COLUMN front_text TO first_side;

ALTER TABLE cards
	RENAME COLUMN back_text TO second_side;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE cards
	RENAME COLUMN first_side TO front_text;

ALTER TABLE cards
	RENAME COLUMN second_side TO back_text;
-- +goose StatementEnd
