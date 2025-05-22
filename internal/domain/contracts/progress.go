package contracts

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
)

type ProgressReader interface {
	GetCardsProgress(ctx context.Context, user entity.UserId, group entity.GroupId) ([]entity.CardProgress, error)
}

type ProgressWriter interface {
	UpdateCardsProgress(
		ctx context.Context,
		user entity.UserId,
		group entity.GroupId,
		cardsProgress []entity.CardProgress,
	) error
}
