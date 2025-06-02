package service

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
)

type UserReader interface {
	GetUser(ctx context.Context, id entity.UserId) (entity.User, error)
}

type UserWriter interface {
	AddUser(ctx context.Context, user entity.User) error
}
