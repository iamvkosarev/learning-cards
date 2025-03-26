package entity

import (
	"errors"
	"google.golang.org/grpc/codes"
)

var (
	ErrGroupExists     = errors.New("group already exists")
	ErrMetadataIsEmpty = errors.New("metadata is empty")
)

type VerificationError struct {
	Message    string
	StatusCode codes.Code
}

func NewVerificationError(message string, statusCode codes.Code) VerificationError {
	return VerificationError{
		Message:    message,
		StatusCode: statusCode,
	}
}

func (e VerificationError) Error() string {
	return e.Message
}
