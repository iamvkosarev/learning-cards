// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: learning_cards/v1/card.proto

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

// Validate checks the field values on ReadingPair with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ReadingPair) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ReadingPair with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ReadingPairMultiError, or
// nil if none found.
func (m *ReadingPair) ValidateAll() error {
	return m.validate(true)
}

func (m *ReadingPair) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Text

	// no validation rules for Reading

	if len(errors) > 0 {
		return ReadingPairMultiError(errors)
	}

	return nil
}

// ReadingPairMultiError is an error wrapping multiple validation errors
// returned by ReadingPair.ValidateAll() if the designated constraints aren't met.
type ReadingPairMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ReadingPairMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ReadingPairMultiError) AllErrors() []error { return m }

// ReadingPairValidationError is the validation error returned by
// ReadingPair.Validate if the designated constraints aren't met.
type ReadingPairValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ReadingPairValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ReadingPairValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ReadingPairValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ReadingPairValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ReadingPairValidationError) ErrorName() string { return "ReadingPairValidationError" }

// Error satisfies the builtin error interface
func (e ReadingPairValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sReadingPair.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ReadingPairValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ReadingPairValidationError{}

// Validate checks the field values on CardSide with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *CardSide) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CardSide with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in CardSideMultiError, or nil
// if none found.
func (m *CardSide) ValidateAll() error {
	return m.validate(true)
}

func (m *CardSide) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetText()) < 1 {
		err := CardSideValidationError{
			field:  "Text",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	for idx, item := range m.GetReadingPairs() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, CardSideValidationError{
						field:  fmt.Sprintf("ReadingPairs[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, CardSideValidationError{
						field:  fmt.Sprintf("ReadingPairs[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return CardSideValidationError{
					field:  fmt.Sprintf("ReadingPairs[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return CardSideMultiError(errors)
	}

	return nil
}

// CardSideMultiError is an error wrapping multiple validation errors returned
// by CardSide.ValidateAll() if the designated constraints aren't met.
type CardSideMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CardSideMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CardSideMultiError) AllErrors() []error { return m }

// CardSideValidationError is the validation error returned by
// CardSide.Validate if the designated constraints aren't met.
type CardSideValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CardSideValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CardSideValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CardSideValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CardSideValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CardSideValidationError) ErrorName() string { return "CardSideValidationError" }

// Error satisfies the builtin error interface
func (e CardSideValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCardSide.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CardSideValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CardSideValidationError{}

// Validate checks the field values on Card with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *Card) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Card with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in CardMultiError, or nil if none found.
func (m *Card) ValidateAll() error {
	return m.validate(true)
}

func (m *Card) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for GroupId

	// no validation rules for FrontText

	// no validation rules for BackText

	// no validation rules for CreatedAt

	for idx, item := range m.GetSides() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, CardValidationError{
						field:  fmt.Sprintf("Sides[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, CardValidationError{
						field:  fmt.Sprintf("Sides[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return CardValidationError{
					field:  fmt.Sprintf("Sides[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return CardMultiError(errors)
	}

	return nil
}

// CardMultiError is an error wrapping multiple validation errors returned by
// Card.ValidateAll() if the designated constraints aren't met.
type CardMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CardMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CardMultiError) AllErrors() []error { return m }

// CardValidationError is the validation error returned by Card.Validate if the
// designated constraints aren't met.
type CardValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CardValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CardValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CardValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CardValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CardValidationError) ErrorName() string { return "CardValidationError" }

// Error satisfies the builtin error interface
func (e CardValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCard.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CardValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CardValidationError{}

// Validate checks the field values on AddCardRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *AddCardRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on AddCardRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in AddCardRequestMultiError,
// or nil if none found.
func (m *AddCardRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *AddCardRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for GroupId

	if utf8.RuneCountInString(m.GetFrontText()) < 1 {
		err := AddCardRequestValidationError{
			field:  "FrontText",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetBackText()) < 1 {
		err := AddCardRequestValidationError{
			field:  "BackText",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(m.GetSidesText()) != 2 {
		err := AddCardRequestValidationError{
			field:  "SidesText",
			reason: "value must contain exactly 2 item(s)",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	for idx, item := range m.GetSidesText() {
		_, _ = idx, item

		if utf8.RuneCountInString(item) < 1 {
			err := AddCardRequestValidationError{
				field:  fmt.Sprintf("SidesText[%v]", idx),
				reason: "value length must be at least 1 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(errors) > 0 {
		return AddCardRequestMultiError(errors)
	}

	return nil
}

// AddCardRequestMultiError is an error wrapping multiple validation errors
// returned by AddCardRequest.ValidateAll() if the designated constraints
// aren't met.
type AddCardRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AddCardRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AddCardRequestMultiError) AllErrors() []error { return m }

// AddCardRequestValidationError is the validation error returned by
// AddCardRequest.Validate if the designated constraints aren't met.
type AddCardRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AddCardRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AddCardRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AddCardRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AddCardRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AddCardRequestValidationError) ErrorName() string { return "AddCardRequestValidationError" }

// Error satisfies the builtin error interface
func (e AddCardRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAddCardRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AddCardRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AddCardRequestValidationError{}

// Validate checks the field values on AddCardResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *AddCardResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on AddCardResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// AddCardResponseMultiError, or nil if none found.
func (m *AddCardResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *AddCardResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for CardId

	if len(errors) > 0 {
		return AddCardResponseMultiError(errors)
	}

	return nil
}

// AddCardResponseMultiError is an error wrapping multiple validation errors
// returned by AddCardResponse.ValidateAll() if the designated constraints
// aren't met.
type AddCardResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AddCardResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AddCardResponseMultiError) AllErrors() []error { return m }

// AddCardResponseValidationError is the validation error returned by
// AddCardResponse.Validate if the designated constraints aren't met.
type AddCardResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AddCardResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AddCardResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AddCardResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AddCardResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AddCardResponseValidationError) ErrorName() string { return "AddCardResponseValidationError" }

// Error satisfies the builtin error interface
func (e AddCardResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAddCardResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AddCardResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AddCardResponseValidationError{}

// Validate checks the field values on GetCardRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *GetCardRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetCardRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in GetCardRequestMultiError,
// or nil if none found.
func (m *GetCardRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetCardRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for CardId

	if len(errors) > 0 {
		return GetCardRequestMultiError(errors)
	}

	return nil
}

// GetCardRequestMultiError is an error wrapping multiple validation errors
// returned by GetCardRequest.ValidateAll() if the designated constraints
// aren't met.
type GetCardRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetCardRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetCardRequestMultiError) AllErrors() []error { return m }

// GetCardRequestValidationError is the validation error returned by
// GetCardRequest.Validate if the designated constraints aren't met.
type GetCardRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetCardRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetCardRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetCardRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetCardRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetCardRequestValidationError) ErrorName() string { return "GetCardRequestValidationError" }

// Error satisfies the builtin error interface
func (e GetCardRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetCardRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetCardRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetCardRequestValidationError{}

// Validate checks the field values on GetCardResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetCardResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetCardResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetCardResponseMultiError, or nil if none found.
func (m *GetCardResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetCardResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetCard()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetCardResponseValidationError{
					field:  "Card",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetCardResponseValidationError{
					field:  "Card",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCard()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetCardResponseValidationError{
				field:  "Card",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return GetCardResponseMultiError(errors)
	}

	return nil
}

// GetCardResponseMultiError is an error wrapping multiple validation errors
// returned by GetCardResponse.ValidateAll() if the designated constraints
// aren't met.
type GetCardResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetCardResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetCardResponseMultiError) AllErrors() []error { return m }

// GetCardResponseValidationError is the validation error returned by
// GetCardResponse.Validate if the designated constraints aren't met.
type GetCardResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetCardResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetCardResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetCardResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetCardResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetCardResponseValidationError) ErrorName() string { return "GetCardResponseValidationError" }

// Error satisfies the builtin error interface
func (e GetCardResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetCardResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetCardResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetCardResponseValidationError{}

// Validate checks the field values on ListCardsRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *ListCardsRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListCardsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListCardsRequestMultiError, or nil if none found.
func (m *ListCardsRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ListCardsRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for GroupId

	if len(errors) > 0 {
		return ListCardsRequestMultiError(errors)
	}

	return nil
}

// ListCardsRequestMultiError is an error wrapping multiple validation errors
// returned by ListCardsRequest.ValidateAll() if the designated constraints
// aren't met.
type ListCardsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListCardsRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListCardsRequestMultiError) AllErrors() []error { return m }

// ListCardsRequestValidationError is the validation error returned by
// ListCardsRequest.Validate if the designated constraints aren't met.
type ListCardsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListCardsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListCardsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListCardsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListCardsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListCardsRequestValidationError) ErrorName() string { return "ListCardsRequestValidationError" }

// Error satisfies the builtin error interface
func (e ListCardsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListCardsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListCardsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListCardsRequestValidationError{}

// Validate checks the field values on ListCardsResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *ListCardsResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ListCardsResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ListCardsResponseMultiError, or nil if none found.
func (m *ListCardsResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ListCardsResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetCards() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ListCardsResponseValidationError{
						field:  fmt.Sprintf("Cards[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ListCardsResponseValidationError{
						field:  fmt.Sprintf("Cards[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ListCardsResponseValidationError{
					field:  fmt.Sprintf("Cards[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ListCardsResponseMultiError(errors)
	}

	return nil
}

// ListCardsResponseMultiError is an error wrapping multiple validation errors
// returned by ListCardsResponse.ValidateAll() if the designated constraints
// aren't met.
type ListCardsResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ListCardsResponseMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ListCardsResponseMultiError) AllErrors() []error { return m }

// ListCardsResponseValidationError is the validation error returned by
// ListCardsResponse.Validate if the designated constraints aren't met.
type ListCardsResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListCardsResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListCardsResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListCardsResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListCardsResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListCardsResponseValidationError) ErrorName() string {
	return "ListCardsResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ListCardsResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListCardsResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListCardsResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListCardsResponseValidationError{}

// Validate checks the field values on UpdateCardRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *UpdateCardRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateCardRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UpdateCardRequestMultiError, or nil if none found.
func (m *UpdateCardRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateCardRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for CardId

	if utf8.RuneCountInString(m.GetFrontText()) < 1 {
		err := UpdateCardRequestValidationError{
			field:  "FrontText",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetBackText()) < 1 {
		err := UpdateCardRequestValidationError{
			field:  "BackText",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(m.GetSidesText()) > 2 {
		err := UpdateCardRequestValidationError{
			field:  "SidesText",
			reason: "value must contain no more than 2 item(s)",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return UpdateCardRequestMultiError(errors)
	}

	return nil
}

// UpdateCardRequestMultiError is an error wrapping multiple validation errors
// returned by UpdateCardRequest.ValidateAll() if the designated constraints
// aren't met.
type UpdateCardRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateCardRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateCardRequestMultiError) AllErrors() []error { return m }

// UpdateCardRequestValidationError is the validation error returned by
// UpdateCardRequest.Validate if the designated constraints aren't met.
type UpdateCardRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateCardRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateCardRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateCardRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateCardRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateCardRequestValidationError) ErrorName() string {
	return "UpdateCardRequestValidationError"
}

// Error satisfies the builtin error interface
func (e UpdateCardRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateCardRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateCardRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateCardRequestValidationError{}

// Validate checks the field values on DeleteCardRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *DeleteCardRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on DeleteCardRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// DeleteCardRequestMultiError, or nil if none found.
func (m *DeleteCardRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *DeleteCardRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for CardId

	if len(errors) > 0 {
		return DeleteCardRequestMultiError(errors)
	}

	return nil
}

// DeleteCardRequestMultiError is an error wrapping multiple validation errors
// returned by DeleteCardRequest.ValidateAll() if the designated constraints
// aren't met.
type DeleteCardRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m DeleteCardRequestMultiError) Error() string {
	msgs := make([]string, 0, len(m))
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m DeleteCardRequestMultiError) AllErrors() []error { return m }

// DeleteCardRequestValidationError is the validation error returned by
// DeleteCardRequest.Validate if the designated constraints aren't met.
type DeleteCardRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e DeleteCardRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e DeleteCardRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e DeleteCardRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e DeleteCardRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e DeleteCardRequestValidationError) ErrorName() string {
	return "DeleteCardRequestValidationError"
}

// Error satisfies the builtin error interface
func (e DeleteCardRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sDeleteCardRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = DeleteCardRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = DeleteCardRequestValidationError{}
