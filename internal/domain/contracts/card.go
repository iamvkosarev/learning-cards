package contracts

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
)

type GroupWriter interface {
	Add(ctx context.Context, group entity.Group) (int64, error)
	Update(ctx context.Context, group entity.Group) error
}
