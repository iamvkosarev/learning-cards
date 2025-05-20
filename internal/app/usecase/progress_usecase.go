package usecase

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/contracts"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
)

type ProgressUseCaseDeps struct {
	ProgressReader contracts.ProgressReader
}

type ProgressUseCase struct {
	ProgressUseCaseDeps
}

func NewProgressUseCase(deps ProgressUseCaseDeps) *ProgressUseCase {
	return &ProgressUseCase{deps}
}

func (p ProgressUseCase) ListGroupsProgress(ctx context.Context) ([]entity.GroupProgress, error) {
	//TODO implement me
	panic("implement me")
}

func (p ProgressUseCase) GetCardsProgress(ctx context.Context, id entity.GroupId) ([]entity.CardProgress, error) {
	//TODO implement me
	panic("implement me")
}
