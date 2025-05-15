-- +goose Up
-- +goose StatementBegin
ALTER TABLE cards ADD CONSTRAINT front_text_min_size CHECK ( length(front_text) > 0 );
ALTER TABLE cards ADD CONSTRAINT back_text_min_size CHECK ( length(back_text) > 0 );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE cards DROP CONSTRAINT IF EXISTS front_text_min_size;
ALTER TABLE cards DROP CONSTRAINT IF EXISTS back_text_min_size;
-- +goose StatementEnd
