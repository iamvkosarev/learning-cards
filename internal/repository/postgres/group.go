package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"time"

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

func (gr *GroupRepository) ListGroups(ctx context.Context, userId entity.UserId) ([]entity.Group, error) {
	const op = "postgres.GroupRepository.ListGroups"

	rows, err := gr.db.Query(
		ctx,
		`SELECT id, user_id, name, description, visibility, created_at, updated_at
		 FROM groups
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
		var visibility uint8
		var description sql.NullString

		err := rows.Scan(
			&group.Id,
			&group.OwnerId,
			&group.Name,
			&description,
			&visibility,
			&group.CreateTime,
			&group.UpdateTime,
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

func (gr *GroupRepository) AddGroup(ctx context.Context, group entity.Group) (entity.GroupId, error) {
	const op = "postgres.GroupRepository.AddGroup"

	var id int64
	err := gr.db.QueryRow(
		ctx,
		`INSERT INTO groups (user_id, name, description, visibility, updated_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		group.OwnerId,
		group.Name,
		group.Description,
		int8(group.Visibility),
		time.Now(),
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

func (gr *GroupRepository) GetGroup(ctx context.Context, groupId entity.GroupId) (entity.Group, error) {
	const op = "postgres.GroupRepository.GetGroup"

	var group entity.Group
	var visibility int8
	var description sql.NullString

	err := gr.db.QueryRow(
		ctx,
		`SELECT id, user_id, name, description, visibility, created_at, updated_at
		 FROM groups
		 WHERE id = $1`,
		groupId,
	).Scan(
		&group.Id,
		&group.OwnerId,
		&group.Name,
		&description,
		&visibility,
		&group.CreateTime,
		&group.UpdateTime,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return group, fmt.Errorf("%s: error no rows: %w", op, entity.ErrGroupNotFound)
		}
		return group, fmt.Errorf("%s: query error: %w", op, err)
	}

	group.Visibility = entity.GroupVisibility(visibility)

	if description.Valid {
		group.Description = description.String
	} else {
		group.Description = ""
	}

	return group, nil
}

func (gr *GroupRepository) UpdateGroup(ctx context.Context, group entity.Group) error {
	const op = "postgres.GroupRepository.UpdateGroup"

	cmdTag, err := gr.db.Exec(
		ctx,
		`UPDATE groups
		 SET name = $1,
		     description = $2,
		     visibility = $3,
		     updated_at = $4
		 WHERE id = $5`,
		group.Name,
		group.Description,
		int8(group.Visibility),
		time.Now(),
		group.Id,
	)

	if err != nil {
		return fmt.Errorf("%s: update error: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, entity.ErrGroupNotFound)
	}

	return nil
}

func (gr *GroupRepository) DeleteGroup(ctx context.Context, groupId entity.GroupId) error {
	const op = "postgres.GroupRepository.DeleteGroup"

	cmdTag, err := gr.db.Exec(
		ctx,
		`DELETE FROM groups WHERE id = $1`,
		groupId,
	)

	if err != nil {
		return fmt.Errorf("%s: delete error: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, entity.ErrGroupNotFound)
	}

	return nil
}
