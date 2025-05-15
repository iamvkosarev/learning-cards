package usecase

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
)

type ReviewUseCase struct {
}

func NewReviewUseCase() *ReviewUseCase {
	return &ReviewUseCase{}
}

func (r ReviewUseCase) GetReviewCards(ctx context.Context, id entity.GroupId) ([]entity.Card, error) {
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
