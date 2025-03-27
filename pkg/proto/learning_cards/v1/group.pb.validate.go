// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: learning_cards/v1/group.proto

package learning_cards

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on CardsGroup with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *CardsGroup) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CardsGroup with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in CardsGroupMultiError, or
// nil if none found.
func (m *CardsGroup) ValidateAll() error {
	return m.validate(true)
}

func (m *CardsGroup) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for OwnerId

	// no validation rules for Name

	// no validation rules for Description

	// no validation rules for CreatedAt

	// no validation rules for Visibility

	if len(errors) > 0 {
		return CardsGroupMultiError(errors)
	}

	return nil
}

// CardsGroupMultiError is an error wrapping multiple validation errors
// returned by CardsGroup.ValidateAll() if the designated constraints aren't met.
type CardsGroupMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CardsGroupMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CardsGroupMultiError) AllErrors() []error { return m }

// CardsGroupValidationError is the validation error returned by
// CardsGroup.Validate if the designated constraints aren't met.
type CardsGroupValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CardsGroupValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CardsGroupValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CardsGroupValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CardsGroupValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CardsGroupValidationError) ErrorName() string { return "CardsGroupValidationError" }

// Error satisfies the builtin error interface
func (e CardsGroupValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCardsGroup.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CardsGroupValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CardsGroupValidationError{}

// Validate checks the field values on CreateCardsGroupRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateCardsGroupRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateCardsGroupRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateCardsGroupRequestMultiError, or nil if none found.
func (m *CreateCardsGroupRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateCardsGroupRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetGroupName()) < 1 {
		err := CreateCardsGroupRequestValidationError{
			field:  "GroupName",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Description

	// no validation rules for Visibility

	if len(errors) > 0 {
		return CreateCardsGroupRequestMultiError(errors)
	}

	return nil
}

// CreateCardsGroupRequestMultiError is an error wrapping multiple validation
// errors returned by CreateCardsGroupRequest.ValidateAll() if the designated
// constraints aren't met.
type CreateCardsGroupRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateCardsGroupRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateCardsGroupRequestMultiError) AllErrors() []error { return m }

// CreateCardsGroupRequestValidationError is the validation error returned by
// CreateCardsGroupRequest.Validate if the designated constraints aren't met.
type CreateCardsGroupRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateCardsGroupRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateCardsGroupRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateCardsGroupRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateCardsGroupRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateCardsGroupRequestValidationError) ErrorName() string {
	return "CreateCardsGroupRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreateCardsGroupRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateCardsGroupRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateCardsGroupRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateCardsGroupRequestValidationError{}

// Validate checks the field values on CreateCardsGroupResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateCardsGroupResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateCardsGroupResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateCardsGroupResponseMultiError, or nil if none found.
func (m *CreateCardsGroupResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateCardsGroupResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for GroupId

	if len(errors) > 0 {
		return CreateCardsGroupResponseMultiError(errors)
	}

	return nil
}

// CreateCardsGroupResponseMultiError is an error wrapping multiple validation
// errors returned by CreateCardsGroupResponse.ValidateAll() if the designated
// constraints aren't met.
type CreateCardsGroupResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateCardsGroupResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateCardsGroupResponseMultiError) AllErrors() []error { return m }

// CreateCardsGroupResponseValidationError is the validation error returned by
// CreateCardsGroupResponse.Validate if the designated constraints aren't met.
type CreateCardsGroupResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateCardsGroupResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateCardsGroupResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateCardsGroupResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateCardsGroupResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateCardsGroupResponseValidationError) ErrorName() string {
	return "CreateCardsGroupResponseValidationError"
}

// Error satisfies the builtin error interface
func (e CreateCardsGroupResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateCardsGroupResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateCardsGroupResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateCardsGroupResponseValidationError{}

// Validate checks the field values on GetCardsGroupRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetCardsGroupRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetCardsGroupRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetCardsGroupRequestMultiError, or nil if none found.
func (m *GetCardsGroupRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetCardsGroupRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for GroupId

	if len(errors) > 0 {
		return GetCardsGroupRequestMultiError(errors)
	}

	return nil
}

// GetCardsGroupRequestMultiError is an error wrapping multiple validation
// errors returned by GetCardsGroupRequest.ValidateAll() if the designated
// constraints aren't met.
type GetCardsGroupRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetCardsGroupRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetCardsGroupRequestMultiError) AllErrors() []error { return m }

// GetCardsGroupRequestValidationError is the validation error returned by
// GetCardsGroupRequest.Validate if the designated constraints aren't met.
type GetCardsGroupRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetCardsGroupRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetCardsGroupRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetCardsGroupRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetCardsGroupRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetCardsGroupRequestValidationError) ErrorName() string {
	return "GetCardsGroupRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetCardsGroupRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetCardsGroupRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetCardsGroupRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetCardsGroupRequestValidationError{}

// Validate checks the field values on GetCardsGroupResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetCardsGroupResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetCardsGroupResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetCardsGroupResponseMultiError, or nil if none found.
func (m *GetCardsGroupResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetCardsGroupResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetGroup()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetCardsGroupResponseValidationError{
					field:  "Group",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetCardsGroupResponseValidationError{
					field:  "Group",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetGroup()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetCardsGroupResponseValidationError{
				field:  "Group",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return GetCardsGroupResponseMultiError(errors)
	}

	return nil
}

// GetCardsGroupResponseMultiError is an error wrapping multiple validation
// errors returned by GetCardsGroupResponse.ValidateAll() if the designated
// constraints aren't met.
type GetCardsGroupResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetCardsGroupResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetCardsGroupResponseMultiError) AllErrors() []error { return m }

// GetCardsGroupResponseValidationError is the validation error returned by
// GetCardsGroupResponse.Validate if the designated constraints aren't met.
type GetCardsGroupResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetCardsGroupResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetCardsGroupResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetCardsGroupResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetCardsGroupResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetCardsGroupResponseValidationError) ErrorName() string {
	return "GetCardsGroupResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetCardsGroupResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetCardsGroupResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetCardsGroupResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetCardsGroupResponseValidationError{}

// Validate checks the field values on ListCardsGroupsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListCardsGroupsResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListCardsGroupsResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListCardsGroupsResponseMultiError, or nil if none found.
func (m *ListCardsGroupsResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ListCardsGroupsResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetGroups() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ListCardsGroupsResponseValidationError{
						field:  fmt.Sprintf("Groups[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ListCardsGroupsResponseValidationError{
						field:  fmt.Sprintf("Groups[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListCardsGroupsResponseValidationError{
					field:  fmt.Sprintf("Groups[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ListCardsGroupsResponseMultiError(errors)
	}

	return nil
}

// ListCardsGroupsResponseMultiError is an error wrapping multiple validation
// errors returned by ListCardsGroupsResponse.ValidateAll() if the designated
// constraints aren't met.
type ListCardsGroupsResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListCardsGroupsResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListCardsGroupsResponseMultiError) AllErrors() []error { return m }

// ListCardsGroupsResponseValidationError is the validation error returned by
// ListCardsGroupsResponse.Validate if the designated constraints aren't met.
type ListCardsGroupsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListCardsGroupsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListCardsGroupsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListCardsGroupsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListCardsGroupsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListCardsGroupsResponseValidationError) ErrorName() string {
	return "ListCardsGroupsResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListCardsGroupsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListCardsGroupsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListCardsGroupsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListCardsGroupsResponseValidationError{}

// Validate checks the field values on UpdateCardsGroupRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateCardsGroupRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateCardsGroupRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdateCardsGroupRequestMultiError, or nil if none found.
func (m *UpdateCardsGroupRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateCardsGroupRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for GroupId

	// no validation rules for GroupName

	// no validation rules for Description

	// no validation rules for Visibility

	if len(errors) > 0 {
		return UpdateCardsGroupRequestMultiError(errors)
	}

	return nil
}

// UpdateCardsGroupRequestMultiError is an error wrapping multiple validation
// errors returned by UpdateCardsGroupRequest.ValidateAll() if the designated
// constraints aren't met.
type UpdateCardsGroupRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateCardsGroupRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateCardsGroupRequestMultiError) AllErrors() []error { return m }

// UpdateCardsGroupRequestValidationError is the validation error returned by
// UpdateCardsGroupRequest.Validate if the designated constraints aren't met.
type UpdateCardsGroupRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateCardsGroupRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateCardsGroupRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateCardsGroupRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateCardsGroupRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateCardsGroupRequestValidationError) ErrorName() string {
	return "UpdateCardsGroupRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateCardsGroupRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateCardsGroupRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateCardsGroupRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateCardsGroupRequestValidationError{}

// Validate checks the field values on DeleteCardsGroupRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DeleteCardsGroupRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteCardsGroupRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeleteCardsGroupRequestMultiError, or nil if none found.
func (m *DeleteCardsGroupRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteCardsGroupRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for GroupId

	if len(errors) > 0 {
		return DeleteCardsGroupRequestMultiError(errors)
	}

	return nil
}

// DeleteCardsGroupRequestMultiError is an error wrapping multiple validation
// errors returned by DeleteCardsGroupRequest.ValidateAll() if the designated
// constraints aren't met.
type DeleteCardsGroupRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteCardsGroupRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteCardsGroupRequestMultiError) AllErrors() []error { return m }

// DeleteCardsGroupRequestValidationError is the validation error returned by
// DeleteCardsGroupRequest.Validate if the designated constraints aren't met.
type DeleteCardsGroupRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteCardsGroupRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteCardsGroupRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteCardsGroupRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteCardsGroupRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteCardsGroupRequestValidationError) ErrorName() string {
	return "DeleteCardsGroupRequestValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteCardsGroupRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteCardsGroupRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteCardsGroupRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteCardsGroupRequestValidationError{}
