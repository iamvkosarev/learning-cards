package contracts

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
)

type GroupReader interface {
	GetGroup(ctx context.Context, groupId entity.GroupId) (entity.Group, error)
	ListGroups(ctx context.Context, id entity.UserId) ([]entity.Group, error)
}

type GroupWriter interface {
	AddGroup(ctx context.Context, group entity.Group) (entity.GroupId, error)
	UpdateGroup(ctx context.Context, group entity.Group) error
	DeleteGroup(ctx context.Context, groupId entity.GroupId) error
}
