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

// Validate checks the field values on CreateGroupRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateGroupRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateGroupRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateGroupRequestMultiError, or nil if none found.
func (m *CreateGroupRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateGroupRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetGroupName()) < 1 {
		err := CreateGroupRequestValidationError{
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
		return CreateGroupRequestMultiError(errors)
	}

	return nil
}

// CreateGroupRequestMultiError is an error wrapping multiple validation errors
// returned by CreateGroupRequest.ValidateAll() if the designated constraints
// aren't met.
type CreateGroupRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateGroupRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateGroupRequestMultiError) AllErrors() []error { return m }

// CreateGroupRequestValidationError is the validation error returned by
// CreateGroupRequest.Validate if the designated constraints aren't met.
type CreateGroupRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateGroupRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateGroupRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateGroupRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateGroupRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateGroupRequestValidationError) ErrorName() string {
	return "CreateGroupRequestValidationError"
}

// Error satisfies the builtin error interface
func (e CreateGroupRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateGroupRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateGroupRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateGroupRequestValidationError{}

// Validate checks the field values on CreateGroupResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CreateGroupResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateGroupResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CreateGroupResponseMultiError, or nil if none found.
func (m *CreateGroupResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateGroupResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for GroupId

	if len(errors) > 0 {
		return CreateGroupResponseMultiError(errors)
	}

	return nil
}

// CreateGroupResponseMultiError is an error wrapping multiple validation
// errors returned by CreateGroupResponse.ValidateAll() if the designated
// constraints aren't met.
type CreateGroupResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateGroupResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateGroupResponseMultiError) AllErrors() []error { return m }

// CreateGroupResponseValidationError is the validation error returned by
// CreateGroupResponse.Validate if the designated constraints aren't met.
type CreateGroupResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateGroupResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateGroupResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateGroupResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateGroupResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateGroupResponseValidationError) ErrorName() string {
	return "CreateGroupResponseValidationError"
}

// Error satisfies the builtin error interface
func (e CreateGroupResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateGroupResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateGroupResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateGroupResponseValidationError{}

// Validate checks the field values on GetGroupRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetGroupRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetGroupRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetGroupRequestMultiError, or nil if none found.
func (m *GetGroupRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetGroupRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for GroupId

	if len(errors) > 0 {
		return GetGroupRequestMultiError(errors)
	}

	return nil
}

// GetGroupRequestMultiError is an error wrapping multiple validation errors
// returned by GetGroupRequest.ValidateAll() if the designated constraints
// aren't met.
type GetGroupRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetGroupRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetGroupRequestMultiError) AllErrors() []error { return m }

// GetGroupRequestValidationError is the validation error returned by
// GetGroupRequest.Validate if the designated constraints aren't met.
type GetGroupRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetGroupRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetGroupRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetGroupRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetGroupRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetGroupRequestValidationError) ErrorName() string { return "GetGroupRequestValidationError" }

// Error satisfies the builtin error interface
func (e GetGroupRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetGroupRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetGroupRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetGroupRequestValidationError{}

// Validate checks the field values on GetGroupResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetGroupResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetGroupResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetGroupResponseMultiError, or nil if none found.
func (m *GetGroupResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetGroupResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetGroup()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetGroupResponseValidationError{
					field:  "Group",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetGroupResponseValidationError{
					field:  "Group",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetGroup()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetGroupResponseValidationError{
				field:  "Group",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return GetGroupResponseMultiError(errors)
	}

	return nil
}

// GetGroupResponseMultiError is an error wrapping multiple validation errors
// returned by GetGroupResponse.ValidateAll() if the designated constraints
// aren't met.
type GetGroupResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetGroupResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetGroupResponseMultiError) AllErrors() []error { return m }

// GetGroupResponseValidationError is the validation error returned by
// GetGroupResponse.Validate if the designated constraints aren't met.
type GetGroupResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetGroupResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetGroupResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetGroupResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetGroupResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetGroupResponseValidationError) ErrorName() string { return "GetGroupResponseValidationError" }

// Error satisfies the builtin error interface
func (e GetGroupResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetGroupResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetGroupResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetGroupResponseValidationError{}

// Validate checks the field values on ListGroupsRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *ListGroupsRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListGroupsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListGroupsRequestMultiError, or nil if none found.
func (m *ListGroupsRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ListGroupsRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return ListGroupsRequestMultiError(errors)
	}

	return nil
}

// ListGroupsRequestMultiError is an error wrapping multiple validation errors
// returned by ListGroupsRequest.ValidateAll() if the designated constraints
// aren't met.
type ListGroupsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListGroupsRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListGroupsRequestMultiError) AllErrors() []error { return m }

// ListGroupsRequestValidationError is the validation error returned by
// ListGroupsRequest.Validate if the designated constraints aren't met.
type ListGroupsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListGroupsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListGroupsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListGroupsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListGroupsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListGroupsRequestValidationError) ErrorName() string {
	return "ListGroupsRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ListGroupsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListGroupsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListGroupsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListGroupsRequestValidationError{}

// Validate checks the field values on ListGroupsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ListGroupsResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListGroupsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListGroupsResponseMultiError, or nil if none found.
func (m *ListGroupsResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ListGroupsResponse) validate(all bool) error {
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
					errors = append(errors, ListGroupsResponseValidationError{
						field:  fmt.Sprintf("Groups[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ListGroupsResponseValidationError{
						field:  fmt.Sprintf("Groups[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListGroupsResponseValidationError{
					field:  fmt.Sprintf("Groups[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ListGroupsResponseMultiError(errors)
	}

	return nil
}

// ListGroupsResponseMultiError is an error wrapping multiple validation errors
// returned by ListGroupsResponse.ValidateAll() if the designated constraints
// aren't met.
type ListGroupsResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListGroupsResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListGroupsResponseMultiError) AllErrors() []error { return m }

// ListGroupsResponseValidationError is the validation error returned by
// ListGroupsResponse.Validate if the designated constraints aren't met.
type ListGroupsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListGroupsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListGroupsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListGroupsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListGroupsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListGroupsResponseValidationError) ErrorName() string {
	return "ListGroupsResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListGroupsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListGroupsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListGroupsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListGroupsResponseValidationError{}

// Validate checks the field values on UpdateGroupRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *UpdateGroupRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateGroupRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdateGroupRequestMultiError, or nil if none found.
func (m *UpdateGroupRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateGroupRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for GroupId

	// no validation rules for GroupName

	// no validation rules for Description

	// no validation rules for Visibility

	if len(errors) > 0 {
		return UpdateGroupRequestMultiError(errors)
	}

	return nil
}

// UpdateGroupRequestMultiError is an error wrapping multiple validation errors
// returned by UpdateGroupRequest.ValidateAll() if the designated constraints
// aren't met.
type UpdateGroupRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateGroupRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateGroupRequestMultiError) AllErrors() []error { return m }

// UpdateGroupRequestValidationError is the validation error returned by
// UpdateGroupRequest.Validate if the designated constraints aren't met.
type UpdateGroupRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateGroupRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateGroupRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateGroupRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateGroupRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateGroupRequestValidationError) ErrorName() string {
	return "UpdateGroupRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateGroupRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateGroupRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateGroupRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateGroupRequestValidationError{}

// Validate checks the field values on DeleteGroupRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *DeleteGroupRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteGroupRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeleteGroupRequestMultiError, or nil if none found.
func (m *DeleteGroupRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteGroupRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for GroupId

	if len(errors) > 0 {
		return DeleteGroupRequestMultiError(errors)
	}

	return nil
}

// DeleteGroupRequestMultiError is an error wrapping multiple validation errors
// returned by DeleteGroupRequest.ValidateAll() if the designated constraints
// aren't met.
type DeleteGroupRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteGroupRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteGroupRequestMultiError) AllErrors() []error { return m }

// DeleteGroupRequestValidationError is the validation error returned by
// DeleteGroupRequest.Validate if the designated constraints aren't met.
type DeleteGroupRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteGroupRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteGroupRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteGroupRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteGroupRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteGroupRequestValidationError) ErrorName() string {
	return "DeleteGroupRequestValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteGroupRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteGroupRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteGroupRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteGroupRequestValidationError{}
