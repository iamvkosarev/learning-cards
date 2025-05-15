-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS group_reviews (
	user_id BIGINT NOT NULL,
	group_id BIGINT NOT NULL,
	last_review_time TIMESTAMP NOT NULL DEFAULT current_timestamp,
	cards_count numeric(3) NOT NULL CHECK(cards_count <= 100) DEFAULT 20,
	PRIMARY KEY (user_id, group_id),
	FOREIGN KEY(user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
	FOREIGN KEY(group_id) REFERENCES groups(id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS card_reviews (
	user_id BIGINT NOT NULL,
	group_id BIGINT NOT NULL,
	card_id BIGINT NOT NULL,
	last_review_time TIMESTAMP NOT NULL DEFAULT current_timestamp,
	fail_count SMALLINT NOT NULL CHECK(fail_count >= 0),
	hard_count SMALLINT NOT NULL CHECK(hard_count >= 0),
	good_count SMALLINT NOT NULL CHECK(good_count >= 0),
	easy_count SMALLINT NOT NULL CHECK(easy_count >= 0),
	avg_review_time NUMERIC(5, 1) NOT NULL,
	PRIMARY KEY(user_id, group_id, card_id),
	FOREIGN KEY(user_id) REFERENCES users(id) ON UPDATE CASCADE ON DELETE CASCADE,
	FOREIGN KEY(group_id) REFERENCES groups(id) ON UPDATE CASCADE ON DELETE CASCADE,
	FOREIGN KEY(card_id) REFERENCES cards(id) ON UPDATE CASCADE ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS group_reviews;
DROP TABLE IF EXISTS card_progress;
-- +goose StatementEnd
