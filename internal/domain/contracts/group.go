package contracts

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
)

type GroupReader interface {
	Get(ctx context.Context, groupId entity.GroupId) (entity.Group, error)
	ListByUser(ctx context.Context, id entity.UserId) ([]entity.Group, error)
}

type GroupWriter interface {
	Add(ctx context.Context, group entity.Group) (entity.GroupId, error)
	Update(ctx context.Context, group entity.Group) error
	Delete(ctx context.Context, groupId entity.GroupId) error
}
