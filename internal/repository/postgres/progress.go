package postgres

import (
	"context"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type ProgressRepository struct {
	db *pgxpool.Pool
}

func NewProgressRepository(pool *pgxpool.Pool) *ProgressRepository {
	return &ProgressRepository{db: pool}
}

func (p ProgressRepository) GetCardsProgress(
	ctx context.Context,
	user entity.UserId,
	group entity.GroupId,
) ([]entity.CardProgress, error) {
	op := "postgres.ProgressRepository.GetCardsMarks"

	rows, err := p.db.Query(
		ctx,
		`SELECT card_id, last_review_time, fail_count, hard_count, good_count, easy_count, avg_review_time FROM card_reviews WHERE user_id = $1 AND group_id = $2`,
		user, group,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: query error: %w", op, err)
	}
	defer rows.Close()
	var cards []entity.CardProgress
	var avgReviewTime float64
	for rows.Next() {
		var card entity.CardProgress
		err = rows.Scan(
			&card.Id,
			&card.LastReviewTime,
			&card.FailsCount,
			&card.HardCount,
			&card.GoodCount,
			&card.EasyCount,
			&avgReviewTime,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: row scan: %w", op, err)
		}
		card.AverageReviewTime = time.Duration(avgReviewTime * float64(time.Second))
		cards = append(cards, card)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error %w", op, err)
	}
	return cards, nil
}

func (p ProgressRepository) UpdateCardsProgress(
	ctx context.Context,
	user entity.UserId,
	group entity.GroupId,
	cardsProgress []entity.CardProgress,
) error {
	op := "postgres.ProgressRepository.UpdateCardsProgress"

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: begin tx: %w", op, err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	stmt := `
	INSERT INTO card_reviews (
		user_id, group_id, card_id,
		last_review_time, fail_count, hard_count, good_count, easy_count, avg_review_time
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	ON CONFLICT (user_id, group_id, card_id) DO UPDATE
	SET
		last_review_time = EXCLUDED.last_review_time,
		fail_count = EXCLUDED.fail_count,
		hard_count = EXCLUDED.hard_count,
		good_count = EXCLUDED.good_count,
		easy_count = EXCLUDED.easy_count,
		avg_review_time = EXCLUDED.avg_review_time
	`

	for _, card := range cardsProgress {
		_, err = tx.Exec(
			ctx, stmt,
			user, group, card.Id,
			card.LastReviewTime,
			card.FailsCount, card.HardCount, card.GoodCount, card.EasyCount,
			card.AverageReviewTime.Seconds(),
		)
		if err != nil {
			return fmt.Errorf("%s: exec for card_id %d: %w", op, card.Id, err)
		}
	}

	return nil
}
