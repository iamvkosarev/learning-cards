package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
	"github.com/iamvkosarev/learning-cards/internal/infrastructure/repo"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProgressRepository struct {
	db *pgxpool.Pool
}

func NewProgressRepository(pool *pgxpool.Pool) *ProgressRepository {
	return &ProgressRepository{db: pool}
}

func (p ProgressRepository) GetGroupProgress(
	ctx context.Context,
	user entity.UserId,
	group entity.GroupId,
) (entity.GroupProgress, error) {
	const op = "postgres.ProgressRepository.GetGroupsProgress"

	groupProgress := entity.GroupProgress{UserId: user, GroupId: group}
	err := p.db.QueryRow(
		ctx, `SELECT last_review_time FROM group_progress WHERE user_id = $1 AND group_id = $2`, user,
		group,
	).Scan(&groupProgress.LastReviewTime)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return groupProgress, fmt.Errorf("%s: %w", op, repository.ErrGroupProgressNotFound)
		}
	}
	return groupProgress, nil
}

func (p ProgressRepository) GetCardsProgress(
	ctx context.Context,
	user entity.UserId,
	group entity.GroupId,
) ([]entity.CardProgress, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProgressRepository) UpdateGroupProgress(
	ctx context.Context,
	user entity.UserId,
	groupProgress entity.GroupProgress,
) error {
	//TODO implement me
	panic("implement me")
}

func (p ProgressRepository) UpdateCardsProgress(
	ctx context.Context,
	user entity.UserId,
	cardsProgress []entity.CardProgress,
) error {
	//TODO implement me
	panic("implement me")
}
