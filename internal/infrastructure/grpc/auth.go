package grpc

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
)

type AuthVerifier interface {
	VerifyUserByContext(ctx context.Context) (userID entity.UserId, err error)
}
