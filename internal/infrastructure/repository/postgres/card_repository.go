package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CardRepository struct {
	db *pgxpool.Pool
}

func NewCardRepository(pool *pgxpool.Pool) *CardRepository {
	return &CardRepository{db: pool}
}

func (cr CardRepository) AddCard(ctx context.Context, card entity.Card) (entity.CardId, error) {
	const op = "postgres.CardRepository.AddCard"

	var id int64
	err := cr.db.QueryRow(
		ctx,
		`INSERT INTO cards (group_id, front_text, back_text) VALUES ($1, $2, $3) RETURNING id`,
		card.GroupId, card.FrontText, card.BackText,
	).Scan(&id)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return 0, fmt.Errorf("%s: %w", op, pgErr)
	}

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return entity.CardId(id), nil
}

func (cr CardRepository) GetCard(ctx context.Context, cardId entity.CardId) (entity.Card, error) {
	op := "postgres.CardRepository.GetCard"

	var card entity.Card

	err := cr.db.QueryRow(
		ctx,
		`SELECT id, group_id, front_text, back_text, created_at, update_at FROM cards WHERE id = $1`,
		cardId,
	).Scan(
		&card.Id,
		&card.GroupId,
		&card.FrontText,
		&card.BackText,
		&card.CreateTime,
		&card.UpdateTime,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return card, fmt.Errorf("%s: card not found: %w", op, entity.ErrCardNotFound)
		}
		return card, fmt.Errorf("%s: %w", op, err)
	}

	return card, nil
}

func (cr CardRepository) ListCards(ctx context.Context, groupId entity.GroupId) ([]entity.Card, error) {
	op := "postgres.CardRepository.ListCards"

	rows, err := cr.db.Query(
		ctx,
		`SELECT id, group_id, front_text, back_text, created_at, update_at
		FROM cards WHERE group_id = $1`, groupId,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: query error: %w", op, err)
	}
	defer rows.Close()

	var cards []entity.Card

	for rows.Next() {
		var card entity.Card
		err = rows.Scan(
			&card.Id,
			&card.GroupId,
			&card.FrontText,
			&card.BackText,
			&card.CreateTime,
			&card.UpdateTime,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan error %w", op, err)
		}
		cards = append(cards, card)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error %w", op, err)
	}

	return cards, nil
}

func (cr CardRepository) UpdateCard(ctx context.Context, card entity.Card) error {
	op := "postgres.CardRepository.UpdateCard"

	cmdTag, err := cr.db.Exec(
		ctx,
		`UPDATE cards
		 SET front_text = $1,
		     back_text = $2,
		     update_at = $3
		 WHERE id = $4`,
		card.FrontText,
		card.BackText,
		card.UpdateTime,
		card.Id,
	)

	if err != nil {
		return fmt.Errorf("%s: update error: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, entity.ErrCardNotFound)
	}

	return nil
}

func (cr CardRepository) DeleteCard(ctx context.Context, cardId entity.CardId) error {
	const op = "postgres.CardRepository.DeleteCard"

	cmdTag, err := cr.db.Exec(
		ctx,
		`DELETE FROM cards WHERE id = $1`,
		cardId,
	)

	if err != nil {
		return fmt.Errorf("%s: delete error: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, entity.ErrCardNotFound)
	}

	return nil
}
