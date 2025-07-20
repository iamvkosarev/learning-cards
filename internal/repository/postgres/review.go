package postgres

import (
	"context"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"slices"
	"strings"
	"time"
)

const reviewTracerName = "postgres.CardRepository"

type ReviewRepository struct {
	db     *pgxpool.Pool
	tracer trace.Tracer
}

func NewReviewRepository(pool *pgxpool.Pool) *ReviewRepository {
	return &ReviewRepository{db: pool, tracer: otel.Tracer(reviewTracerName)}
}

func (p *ReviewRepository) DeleteNotUsedReviews(
	ctx context.Context,
	userId model.UserId,
	groupId model.GroupId,
) error {
	ctx, span := p.tracer.Start(ctx, "DeleteNotUsedReviews")
	defer span.End()

	op := "postgres.ReviewRepository.DeleteNotUsedReviews"

	reviews, err := p.GetCardsReviews(ctx, userId, groupId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	cardReviewsMap := make(map[model.CardId][]*model.CardReview)
	for _, review := range reviews {
		cardReviewsMap[review.CardId] = append(cardReviewsMap[review.CardId], review)
	}

	var toDelete []struct {
		CardId model.CardId
		Time   time.Time
	}

	for _, cardReviews := range cardReviewsMap {
		if len(cardReviews) <= model.MaxReviewsPerCard {
			continue
		}

		slices.SortFunc(
			cardReviews, func(a, b *model.CardReview) int {
				return b.Time.Compare(a.Time)
			},
		)

		for _, r := range cardReviews[model.MaxReviewsPerCard:] {
			toDelete = append(
				toDelete, struct {
					CardId model.CardId
					Time   time.Time
				}{
					CardId: r.CardId,
					Time:   r.Time,
				},
			)
		}
	}

	if len(toDelete) == 0 {
		return nil
	}

	sb := strings.Builder{}
	args := []interface{}{userId, groupId}
	sb.WriteString(
		`
		DELETE FROM card_reviews
		WHERE user_id = $1 AND group_id = $2 AND (card_id, time) IN (`,
	)

	paramIdx := 3
	for i, r := range toDelete {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("($%d, $%d)", paramIdx, paramIdx+1))
		args = append(args, r.CardId, r.Time)
		paramIdx += 2
	}
	sb.WriteString(")")

	_, err = p.db.Exec(ctx, sb.String(), args...)
	if err != nil {
		return fmt.Errorf("%s: delete query: %w", op, err)
	}

	return nil
}

func (p *ReviewRepository) GetCardsReviews(
	ctx context.Context,
	user model.UserId,
	group model.GroupId,
) ([]*model.CardReview, error) {
	ctx, span := p.tracer.Start(ctx, "GetCardsReviews")
	defer span.End()

	op := "postgres.ReviewRepository.GetCardsMarks"

	rows, err := p.db.Query(
		ctx,
		`SELECT card_id, time, duration, answer FROM card_reviews WHERE user_id = $1 AND group_id = $2`,
		user, group,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: query error: %w", op, err)
	}
	defer rows.Close()
	cards := make([]*model.CardReview, 0)
	var duration float64
	for rows.Next() {
		review := &model.CardReview{}
		err = rows.Scan(
			&review.CardId,
			&review.Time,
			&duration,
			&review.Answer,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: row scan: %w", op, err)
		}
		review.Duration = time.Duration(duration * float64(time.Second))
		review.UserId = user
		review.GroupId = group
		cards = append(cards, review)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error %w", op, err)
	}
	return cards, nil
}

func (p *ReviewRepository) AddCardsReviews(
	ctx context.Context,
	user model.UserId,
	group model.GroupId,
	cardsProgress []model.CardReview,
) error {
	ctx, span := p.tracer.Start(ctx, "AddCardsReviews")
	defer span.End()

	op := "postgres.ReviewRepository.AddCardsReviews"

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
		time, duration, answer
	)
	VALUES ($1, $2, $3, $4, $5, $6)`

	for _, card := range cardsProgress {
		_, err = tx.Exec(
			ctx, stmt,
			user, group, card.CardId,
			card.Time, card.Duration.Seconds(), card.Answer,
		)
		if err != nil {
			return fmt.Errorf("%s: exec for card_id %d: %w", op, card.CardId, err)
		}
	}

	return nil
}
