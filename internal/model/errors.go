package model

import (
	"google.golang.org/grpc/codes"
)

var (
	ErrMetadataIsEmpty = NewServerError(codes.InvalidArgument, "metadata is empty")

	ErrNoAuthHeader        = NewServerError(codes.InvalidArgument, "there is no authorization header")
	ErrIncorrectAuthHeader = NewServerError(codes.InvalidArgument, "not correct authorization header")
	ErrVerificationFailed  = NewServerError(codes.PermissionDenied, "verification failed")

	ErrUserNotFound = NewServerError(codes.NotFound, "user not found")

	ErrGroupExists                     = NewServerError(codes.InvalidArgument, "group already exists")
	ErrGroupNotFound                   = NewServerError(codes.NotFound, "group not found")
	ErrGroupReadAccessDenied           = NewServerError(codes.PermissionDenied, "group read access denied")
	ErrGroupWriteAccessDenied          = NewServerError(codes.PermissionDenied, "group write access denied")
	ErrGroupModifyNotNullCardsSideType = NewServerError(
		codes.InvalidArgument,
		"forbidden change of group's not-null cards-side-type value",
	)

	ErrCardNotFound = NewServerError(codes.NotFound, "card not found")

	ErrTimeOut = NewServerError(codes.Canceled, "time is out")
)

type ValidationError struct {
	err error
}

func NewValidationError(err error) *ValidationError {
	return &ValidationError{err: err}
}

func (v *ValidationError) Error() string {
	return v.err.Error()
}

type ServerError struct {
	Code    codes.Code
	Message string
}

func NewServerError(c codes.Code, msg string) error {
	return &ServerError{
		Code:    c,
		Message: msg,
	}
}

func (s *ServerError) Error() string {
	return s.Message
}
