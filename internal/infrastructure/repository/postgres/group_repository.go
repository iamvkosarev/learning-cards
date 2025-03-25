package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
)

type GroupRepository struct {
	db *pgxpool.Pool
}

func NewGroupRepository(pool *pgxpool.Pool) *GroupRepository {
	return &GroupRepository{db: pool}
}

func (r *GroupRepository) Add(ctx context.Context, user entity.CardGroup) (int64, error) {
	const op = "repository.postgres.Add"

	var id int64
	err := r.db.QueryRow(
		ctx,
		`INSERT INTO groups (name, user_id) VALUES ($1, $2) RETURNING id`,
		user.Name, user.UserID,
	).Scan(&id)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, entity.ErrGroupExists)
		}
		return 0, fmt.Errorf("%s: %w", op, pgErr)
	}

	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
