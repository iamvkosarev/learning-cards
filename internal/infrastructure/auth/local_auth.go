package auth

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/domain/entity"
)

type LocalAuth struct {
	userId int
}

func NewLocalService(userId int) *LocalAuth {
	return &LocalAuth{userId: userId}
}

func (l *LocalAuth) VerifyUserByContext(_ context.Context) (userID entity.UserId, err error) {
	return entity.UserId(l.userId), nil
}
