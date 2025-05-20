package usecase

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/contracts"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
)

type ReviewUseCaseDeps struct {
	ProgressReader contracts.ProgressReader
	ProgressWriter contracts.ProgressWriter
}

type ReviewUseCase struct {
	ReviewUseCaseDeps
}

func NewReviewUseCase(deps ReviewUseCaseDeps) *ReviewUseCase {
	return &ReviewUseCase{deps}
}

func (r ReviewUseCase) GetReviewCards(ctx context.Context, id entity.GroupId) ([]entity.Card, error) {
	//TODO implement me
	panic("implement me")
}

func (r ReviewUseCase) SaveResults(ctx context.Context, answers []entity.ReviewCardResult) error {
	//TODO implement me
	panic("implement me")
}

func (r ReviewUseCase) GetGroupReviewInfo(ctx context.Context, id entity.GroupId) (entity.GroupReviewInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (r ReviewUseCase) UpdateGroupReviewInfo(ctx context.Context, reviewInfo entity.UpdateGroupReviewInfo) error {
	//TODO implement me
	panic("implement me")
}
