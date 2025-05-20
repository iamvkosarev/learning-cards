package contracts

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
)

type ProgressReader interface {
	GetGroupProgress(ctx context.Context, user entity.UserId, group entity.GroupId) (entity.GroupProgress, error)
	GetCardsProgress(ctx context.Context, user entity.UserId, group entity.GroupId) ([]entity.CardProgress, error)
}

type ProgressWriter interface {
	UpdateGroupProgress(ctx context.Context, user entity.UserId, groupProgress entity.GroupProgress) error
	UpdateCardsProgress(ctx context.Context, user entity.UserId, cardsProgress []entity.CardProgress) error
}
