package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/model"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GroupRepository struct {
	db *pgxpool.Pool
}

func NewGroupRepository(pool *pgxpool.Pool) *GroupRepository {
	return &GroupRepository{db: pool}
}

func (gr *GroupRepository) ListGroups(ctx context.Context, userId model.UserId) ([]*model.Group, error) {
	const op = "postgres.GroupRepository.ListGroups"

	rows, err := gr.db.Query(
		ctx,
		`SELECT id, user_id, name, description, visibility, created_at, updated_at, first_side_type, second_side_type
		 FROM groups
		 WHERE user_id = $1`,
		userId,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: query error: %w", op, err)
	}
	defer rows.Close()

	groups := make([]*model.Group, 0)

	for rows.Next() {
		var description sql.NullString
		group := &model.Group{
			CardSideTypes: make([]model.CardSideType, 2),
		}

		err = rows.Scan(
			&group.Id,
			&group.OwnerId,
			&group.Name,
			&description,
			&group.Visibility,
			&group.CreateTime,
			&group.UpdateTime,
			&group.CardSideTypes[model.CARD_SIDE_FIRST],
			&group.CardSideTypes[model.CARD_SIDE_SECOND],
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan error: %w", op, err)
		}
		if description.Valid {
			group.Description = description.String
		} else {
			group.Description = ""
		}
		groups = append(groups, group)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, rows.Err())
	}

	return groups, nil
}

func (gr *GroupRepository) AddGroup(ctx context.Context, group *model.Group) (model.GroupId, error) {
	const op = "postgres.GroupRepository.AddGroup"

	var id int64
	err := gr.db.QueryRow(
		ctx,
		`INSERT INTO groups (user_id, name, description, visibility, updated_at, first_side_type, second_side_type)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id`,
		group.OwnerId,
		group.Name,
		group.Description,
		group.Visibility,
		time.Now(),
		group.CardSideTypes[model.CARD_SIDE_FIRST],
		group.CardSideTypes[model.CARD_SIDE_SECOND],
	).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return 0, fmt.Errorf("%s: %w", op, model.ErrGroupExists)
			default:
				return 0, fmt.Errorf("%s: postgres error: %w", op, pgErr)
			}
		}
		return 0, fmt.Errorf("%s: query error: %w", op, err)
	}

	return model.GroupId(id), nil
}

func (gr *GroupRepository) GetGroup(ctx context.Context, groupId model.GroupId) (*model.Group, error) {
	const op = "postgres.GroupRepository.GetGroup"

	group := &model.Group{
		CardSideTypes: make([]model.CardSideType, 2),
	}
	var description sql.NullString

	err := gr.db.QueryRow(
		ctx,
		`SELECT id, user_id, name, description, visibility, created_at, updated_at, first_side_type, second_side_type
		 FROM groups
		 WHERE id = $1`,
		groupId,
	).Scan(
		&group.Id,
		&group.OwnerId,
		&group.Name,
		&description,
		&group.Visibility,
		&group.CreateTime,
		&group.UpdateTime,
		&group.CardSideTypes[model.CARD_SIDE_FIRST],
		&group.CardSideTypes[model.CARD_SIDE_SECOND],
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return group, fmt.Errorf("%s: error no rows: %w", op, model.ErrGroupNotFound)
		}
		return group, fmt.Errorf("%s: query error: %w", op, err)
	}

	if description.Valid {
		group.Description = description.String
	} else {
		group.Description = ""
	}

	return group, nil
}

func (gr *GroupRepository) UpdateGroup(ctx context.Context, group *model.Group) error {
	const op = "postgres.GroupRepository.UpdateGroup"

	cmdTag, err := gr.db.Exec(
		ctx,
		`UPDATE groups
		 SET name = $1,
		     description = $2,
		     visibility = $3,
		     updated_at = $4,
		     first_side_type = $5,
		     second_side_type = $6
		WHERE id = $7`,
		group.Name,
		group.Description,
		int8(group.Visibility),
		time.Now(),
		group.CardSideTypes[model.CARD_SIDE_FIRST],
		group.CardSideTypes[model.CARD_SIDE_SECOND],
		group.Id,
	)

	if err != nil {
		return fmt.Errorf("%s: update error: %w", op, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, model.ErrGroupNotFound)
	}

	return nil
}

func (gr *GroupRepository) DeleteGroup(ctx context.Context, groupId model.GroupId) error {
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
		return fmt.Errorf("%s: %w", op, model.ErrGroupNotFound)
	}

	return nil
}
