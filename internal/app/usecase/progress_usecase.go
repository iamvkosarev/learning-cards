package usecase

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
)

type ProgressUseCase struct {
}

func NewProgressUseCase() *ProgressUseCase {
	return &ProgressUseCase{}
}

func (p ProgressUseCase) ListGroupsProgress(ctx context.Context) ([]entity.GroupProgress, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProgressUseCase) GetCardsProgress(ctx context.Context, id entity.GroupId) ([]entity.CardProgress, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProgressUseCase) UpdateProgress(ctx context.Context, card []entity.ReviewCardResult) error {
	//TODO implement me
	panic("implement me")
}
