package server

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/model"
)

type AuthVerifier interface {
	VerifyUserByContext(ctx context.Context) (userID model.UserId, err error)
}

type VerifyFunc func(ctx context.Context) (model.UserId, error)

func (v VerifyFunc) VerifyUserByContext(ctx context.Context) (userID model.UserId, err error) {
	return v(ctx)
}
