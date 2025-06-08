-- +goose Up
-- +goose StatementBegin
ALTER TABLE groups
	ADD COLUMN first_side_type NUMERIC(1) DEFAULT 0 NOT NULL,
	ADD COLUMN second_side_type NUMERIC(1) DEFAULT 0 NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE groups
	DROP COLUMN first_side_type,
	DROP COLUMN second_side_type;
-- +goose StatementEnd
