-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS group_reviews;
-- +goose StatementEnd

-- +goose Down
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
-- +goose StatementEnd
