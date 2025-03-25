package usecase

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
)

type groupRepository interface {
	Add(ctx context.Context, group entity.CardGroup) (int64, error)
}

type GroupUseCase struct {
	groupRepository
}

func NewGroupUseCase(repo groupRepository) *GroupUseCase {
	return &GroupUseCase{
		groupRepository: repo,
	}
}

func (uc *GroupUseCase) Create(ctx context.Context, name string, userID int64) (int64, error) {
	group := entity.CardGroup{
		Name:   name,
		UserID: userID,
	}

	id, err := uc.groupRepository.Add(ctx, group)
	if err != nil {
		return 0, err
	}

	return id, nil
}
