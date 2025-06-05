-- +goose Up
-- +goose StatementBegin
ALTER TABLE card_reviews
    ADD COLUMN IF NOT EXISTS answer NUMERIC(1) DEFAULT 0 NOT NULL CHECK (answer >= 0 AND answer <= 4);

ALTER TABLE card_reviews
    DROP COLUMN IF EXISTS fail_count,
    DROP COLUMN IF EXISTS hard_count,
    DROP COLUMN IF EXISTS good_count,
    DROP COLUMN IF EXISTS easy_count;

ALTER TABLE card_reviews
    RENAME COLUMN avg_review_time TO duration;

ALTER TABLE card_reviews
    RENAME COLUMN last_review_time TO time;

ALTER TABLE card_reviews
    DROP CONSTRAINT IF EXISTS card_reviews_pkey,
    DROP CONSTRAINT IF EXISTS card_reviews_fail_count_check,
    DROP CONSTRAINT IF EXISTS card_reviews_hard_count_check,
    DROP CONSTRAINT IF EXISTS card_reviews_good_count_check,
    DROP CONSTRAINT IF EXISTS card_reviews_easy_count_check;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE card_reviews
    DROP COLUMN IF EXISTS answer;

ALTER TABLE card_reviews
    ADD COLUMN IF NOT EXISTS fail_count SMALLINT NOT NULL CHECK (fail_count >= 0),
    ADD COLUMN IF NOT EXISTS hard_count SMALLINT NOT NULL CHECK (hard_count >= 0),
    ADD COLUMN IF NOT EXISTS good_count SMALLINT NOT NULL CHECK (good_count >= 0),
    ADD COLUMN IF NOT EXISTS easy_count SMALLINT NOT NULL CHECK (easy_count >= 0);

ALTER TABLE card_reviews
    RENAME COLUMN duration TO avg_review_time;

ALTER TABLE card_reviews
    RENAME COLUMN time TO last_review_time;

ALTER TABLE card_reviews
    ADD CONSTRAINT card_reviews_pkey PRIMARY KEY (user_id, group_id, card_id);
-- +goose StatementEnd
