package verification

import (
	"context"
)

type StubVerifier struct {
	userId int64
}

func (s StubVerifier) VerifyToken(context.Context, string) (int64, error) {
	return s.userId, nil
}

func NewStubVerifier(userId int64) *StubVerifier {
	return &StubVerifier{userId: userId}
}
