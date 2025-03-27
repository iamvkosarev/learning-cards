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

func (cr CardRepository) Add(ctx context.Context, card entity.Card) (entity.CardId, error) {
	const op = "postgres.CardRepository.Add"

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

func (cr CardRepository) Get(ctx context.Context, cardId entity.CardId) (entity.Card, error) {
	op := "postgres.CardRepository.Get"

	var card entity.Card

	err := cr.db.QueryRow(
		ctx,
		`SELECT id, group_id, front_text, back_text, created_at FROM cards WHERE id = $1`,
		cardId,
	).Scan(
		&card.Id,
		&card.GroupId,
		&card.FrontText,
		&card.BackText,
		&card.CreateTime,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return card, fmt.Errorf("%s: card not found: %w", op, entity.ErrCardNotFound)
		}
		return card, fmt.Errorf("%s: %w", op, err)
	}

	return card, nil
}

func (cr CardRepository) List(ctx context.Context, groupId entity.GroupId) ([]entity.Card, error) {
	op := "postgres.CardRepository.List"

	rows, err := cr.db.Query(
		ctx,
		`SELECT id, group_id, front_text, back_text, created_at
		FROM cards WHERE group_id = $1`, groupId,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: query error: %w", op, err)
	}
	defer rows.Close()

	var cards []entity.Card

	for rows.Next() {
		var card entity.Card
		err := rows.Scan(
			&card.Id,
			&card.GroupId,
			&card.FrontText,
			&card.BackText,
			&card.CreateTime,
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

func (cr CardRepository) Delete(ctx context.Context, cardId entity.CardId) error {
	//TODO implement me
	panic("implement me")
}

func (cr CardRepository) Update(ctx context.Context, card entity.Card) error {
	//TODO implement me
	panic("implement me")
}
