package server

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
)

type AuthVerifier interface {
	VerifyUserByContext(ctx context.Context) (userID entity.UserId, err error)
}

type VerifyFunc func(ctx context.Context) (entity.UserId, error)

func (v VerifyFunc) VerifyUserByContext(ctx context.Context) (userID entity.UserId, err error) {
	return v(ctx)
}
