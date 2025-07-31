package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/model"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const groupTracerName = "repository.group"

type Group struct {
	db     *pgxpool.Pool
	tracer trace.Tracer
}

func NewGroupRepository(pool *pgxpool.Pool) *Group {
	return &Group{db: pool, tracer: otel.Tracer(groupTracerName)}
}

func (gr *Group) ListGroups(ctx context.Context, userId model.UserId) ([]*model.Group, error) {
	ctx, span := gr.tracer.Start(ctx, "ListGroups")
	defer span.End()

	const op = "repository.Group.ListGroups"

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
			&group.CardSideTypes[model.SIDE_FIRST],
			&group.CardSideTypes[model.SIDE_SECOND],
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

func (gr *Group) AddGroup(ctx context.Context, group *model.Group) (model.GroupId, error) {
	ctx, span := gr.tracer.Start(ctx, "AddGroup")
	defer span.End()

	const op = "repository.Group.AddGroup"

	var id int64
	err := gr.db.QueryRow(
		ctx,
		`INSERT INTO groups (user_id, name, description, visibility, updated_at, first_side_type, second_side_type)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id`,
		group.OwnerId,
		group.Name,
		group.Description,
		group.Visibility,
		time.Now(),
		group.CardSideTypes[model.SIDE_FIRST],
		group.CardSideTypes[model.SIDE_SECOND],
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

func (gr *Group) GetGroup(ctx context.Context, groupId model.GroupId) (*model.Group, error) {
	ctx, span := gr.tracer.Start(ctx, "GetGroup")
	defer span.End()

	const op = "repository.Group.GetGroup"

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
		&group.CardSideTypes[model.SIDE_FIRST],
		&group.CardSideTypes[model.SIDE_SECOND],
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

func (gr *Group) UpdateGroup(ctx context.Context, group *model.Group) error {
	ctx, span := gr.tracer.Start(ctx, "UpdateGroup")
	defer span.End()

	const op = "repository.Group.UpdateGroup"

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
		group.CardSideTypes[model.SIDE_FIRST],
		group.CardSideTypes[model.SIDE_SECOND],
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

func (gr *Group) DeleteGroup(ctx context.Context, groupId model.GroupId) error {
	ctx, span := gr.tracer.Start(ctx, "DeleteGroup")
	defer span.End()

	const op = "repository.Group.DeleteGroup"

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
