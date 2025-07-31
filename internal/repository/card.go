package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

const cardsTracerName = "repository.card"

type Card struct {
	db     *pgxpool.Pool
	tracer trace.Tracer
}

func NewCardRepository(pool *pgxpool.Pool) *Card {
	return &Card{db: pool, tracer: otel.Tracer(cardsTracerName)}
}

func (cr *Card) AddCard(ctx context.Context, card *model.Card) (model.CardId, error) {
	ctx, span := cr.tracer.Start(ctx, "AddCard")
	defer span.End()

	const op = "repository.Card.AddCard"

	var id int64
	err := cr.db.QueryRow(
		ctx,
		`INSERT INTO cards (group_id, first_side, second_side) VALUES ($1, $2, $3) RETURNING id`,
		card.GroupId, card.Sides[model.SIDE_FIRST].Text, card.Sides[model.SIDE_SECOND].Text,
	).Scan(&id)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return 0, fmt.Errorf("%s: %w", op, pgErr)
	}

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return model.CardId(id), nil
}

func (cr *Card) GetCard(ctx context.Context, cardId model.CardId) (*model.Card, error) {
	ctx, span := cr.tracer.Start(ctx, "GetCard")
	defer span.End()

	op := "repository.Card.GetCard"

	card := &model.Card{
		Sides: make([]model.CardSide, 2),
	}
	err := cr.db.QueryRow(
		ctx,
		`SELECT id, group_id, first_side, second_side, created_at, updated_at FROM cards WHERE id = $1`,
		cardId,
	).Scan(
		&card.Id,
		&card.GroupId,
		&card.Sides[model.SIDE_FIRST].Text,
		&card.Sides[model.SIDE_SECOND].Text,
		&card.CreateTime,
		&card.UpdateTime,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: failed to parse reading pairs %w", op, err)
	}

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: card id %v not found: %w", op, cardId, model.ErrCardNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return card, nil
}

func (cr *Card) ListCards(ctx context.Context, groupId model.GroupId) ([]*model.Card, error) {
	ctx, span := cr.tracer.Start(ctx, "ListCards")
	defer span.End()

	op := "repository.Card.ListCards"

	rows, err := cr.db.Query(
		ctx,
		`SELECT id, group_id, first_side, second_side, created_at, updated_at
		FROM cards WHERE group_id = $1`, groupId,
	)

	if err != nil {
		return nil, fmt.Errorf("%s: query error: %w", op, err)
	}
	defer rows.Close()

	cards := make([]*model.Card, 0)

	for rows.Next() {
		card := &model.Card{
			Sides: make([]model.CardSide, 2),
		}
		err = rows.Scan(
			&card.Id,
			&card.GroupId,
			&card.Sides[model.SIDE_FIRST].Text,
			&card.Sides[model.SIDE_SECOND].Text,
			&card.CreateTime,
			&card.UpdateTime,
		)

		if err != nil {
			return nil, fmt.Errorf("%s: scan error %w", op, err)
		}
		cards = append(cards, card)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error %w", op, err)
	}

	return cards, nil
}

func (cr *Card) UpdateCard(ctx context.Context, card *model.Card) error {
	ctx, span := cr.tracer.Start(ctx, "UpdateCard")
	defer span.End()

	op := "repository.Card.UpdateCard"

	cmdTag, err := cr.db.Exec(
		ctx,
		`UPDATE cards
		 SET first_side = $1,
		     second_side = $2,
		     updated_at = $3
		 WHERE id = $4`,
		card.Sides[model.SIDE_FIRST].Text,
		card.Sides[model.SIDE_SECOND].Text,
		card.UpdateTime,
		card.Id,
	)

	if err != nil {
		return fmt.Errorf("%s: update error: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: not affected rows, while updating card id %v: %w", op, card.Id, model.ErrCardNotFound)
	}

	return nil
}

func (cr *Card) DeleteCard(ctx context.Context, cardId model.CardId) error {
	ctx, span := cr.tracer.Start(ctx, "DeleteCard")
	defer span.End()

	const op = "repository.Card.DeleteCard"

	cmdTag, err := cr.db.Exec(
		ctx,
		`DELETE FROM cards WHERE id = $1`,
		cardId,
	)

	if err != nil {
		return fmt.Errorf("%s: delete error: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, model.ErrCardNotFound)
	}

	return nil
}
