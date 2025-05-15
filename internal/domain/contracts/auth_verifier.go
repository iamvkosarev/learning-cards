package contracts

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
)

type VerifyFunc func(ctx context.Context) (entity.UserId, error)

func (v VerifyFunc) VerifyUserByContext(ctx context.Context) (userID entity.UserId, err error) {
	return v(ctx)
}

type AuthVerifier interface {
	VerifyUserByContext(ctx context.Context) (userID entity.UserId, err error)
}
