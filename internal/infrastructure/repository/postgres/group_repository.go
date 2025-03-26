package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"

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

func (gr *GroupRepository) ListByUser(ctx context.Context, userId entity.UserId) ([]entity.Group, error) {
	const op = "postgres.GroupRepository.List"

	rows, err := gr.db.Query(
		ctx,
		`SELECT id, user_id, name, description, created_at, visibility
		 FROM card_groups
		 WHERE user_id = $1`,
		userId,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: query error: %w", op, err)
	}
	defer rows.Close()

	var groups []entity.Group

	for rows.Next() {
		var group entity.Group
		var visibility int32
		var description sql.NullString

		err := rows.Scan(
			&group.Id,
			&group.OwnerId,
			&group.Name,
			&description,
			&group.CreateTime,
			&visibility,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan error: %w", op, err)
		}
		if description.Valid {
			group.Description = description.String
		} else {
			group.Description = ""
		}

		group.Visibility = entity.GroupVisibility(visibility)
		groups = append(groups, group)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, rows.Err())
	}

	return groups, nil
}

func (gr *GroupRepository) Add(ctx context.Context, group entity.Group) (entity.GroupId, error) {
	const op = "postgres.GroupRepository.Add"

	var id int64
	err := gr.db.QueryRow(
		ctx,
		`INSERT INTO card_groups (name, user_id, description, visibility)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		group.Name,
		group.OwnerId,
		group.Description,
		int32(group.Visibility),
	).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return 0, fmt.Errorf("%s: %w", op, entity.ErrGroupExists)
			default:
				return 0, fmt.Errorf("%s: postgres error: %w", op, pgErr)
			}
		}
		return 0, fmt.Errorf("%s: query error: %w", op, err)
	}

	return entity.GroupId(id), nil
}

func (gr *GroupRepository) Get(ctx context.Context, groupId entity.GroupId) (entity.Group, error) {
	const op = "postgres.GroupRepository.Get"

	var g entity.Group
	var visibility int32
	var description sql.NullString

	err := gr.db.QueryRow(
		ctx,
		`SELECT id, user_id, name, description, created_at, visibility
		 FROM card_groups
		 WHERE id = $1`,
		groupId,
	).Scan(
		&g.Id,
		&g.OwnerId,
		&g.Name,
		&description,
		&g.CreateTime,
		&visibility,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return g, fmt.Errorf("%s: group not found: %w", op, entity.ErrGroupNotFound)
		}
		return g, fmt.Errorf("%s: query error: %w", op, err)
	}

	g.Visibility = entity.GroupVisibility(visibility)

	if description.Valid {
		g.Description = description.String
	} else {
		g.Description = ""
	}

	return g, nil
}

func (gr *GroupRepository) Update(ctx context.Context, group entity.Group) error {
	//TODO implement me
	panic("implement me")
}

func (gr *GroupRepository) Delete(ctx context.Context, groupId entity.GroupId) error {
	//TODO implement me
	panic("implement me")
}
