package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	db *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *User {
	return &User{db: pool}
}

func (u *User) GetUser(ctx context.Context, id model.UserId) (model.User, error) {
	op := "repository.User.GetUser"

	var user model.User

	err := u.db.QueryRow(
		ctx,
		`SELECT id FROM users WHERE id = $1`,
		id,
	).Scan(
		&user.UserId,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, model.ErrUserNotFound
		}
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *User) AddUser(ctx context.Context, user model.User) error {
	const op = "repository.User.AddUser"

	var id int64
	err := u.db.QueryRow(
		ctx,
		`INSERT INTO users (id) VALUES ($1) RETURNING id`,
		user.UserId,
	).Scan(&id)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return fmt.Errorf("%s: %w", op, pgErr)
	}

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
