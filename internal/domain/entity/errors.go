package entity

import (
	"errors"
	"google.golang.org/grpc/codes"
)

var (
	ErrGroupExists     = errors.New("group already exists")
	ErrGroupNotFound   = errors.New("access denied")
	ErrMetadataIsEmpty = errors.New("metadata is empty")
	ErrCardNotFound    = errors.New("card not found")
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
