package entity

import (
	"errors"
)

var (
	ErrGroupExists         = errors.New("group already exists")
	ErrGroupNotFound       = errors.New("access denied")
	ErrMetadataIsEmpty     = errors.New("metadata is empty")
	ErrNoAuthHeader        = errors.New("there is no authorization header")
	ErrIncorrectAuthHeader = errors.New("not correct authorization header")
	ErrCardNotFound        = errors.New("card not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrVerificationFailed  = errors.New("verification failed")
)

type VerificationError struct {
	err error
}

func NewVerificationError(err error) *VerificationError {
	return &VerificationError{err: err}
}

func (v *VerificationError) Error() string {
	return v.err.Error()
}

type ValidationError struct {
	err error
}

func NewValidationError(err error) *ValidationError {
	return &ValidationError{err: err}
}

func (v *ValidationError) Error() string {
	return v.err.Error()
}
