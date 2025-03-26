package contracts

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
)

type GroupReader interface {
	Get(ctx context.Context, groupId int64) (entity.Group, error)
	List(ctx context.Context, userId int64) ([]entity.Group, error)
}
