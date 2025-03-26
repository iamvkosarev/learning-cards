package usecase

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
)

type authService interface {
	VerifyToken(ctx context.Context) (userID int64, err error)
}

type groupRepository interface {
	Add(ctx context.Context, group entity.CardGroup) (int64, error)
}

type GroupUseCase struct {
	groupRepository
	authService
}

func NewGroupUseCase(repo groupRepository, authService authService) *GroupUseCase {
	return &GroupUseCase{
		groupRepository: repo,
		authService:     authService,
	}
}

func (uc *GroupUseCase) Create(ctx context.Context, name string) (int64, error) {
	userID, err := uc.authService.VerifyToken(ctx)
	if err != nil {
		return 0, err
	}

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
