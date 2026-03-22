package blocks

import (
	"errors"
	"fmt"
)

// Sentinel errors for common validation failures.
var (
	ErrEmptyText         = errors.New("text cannot be empty")
	ErrExceedsMaxLen     = errors.New("exceeds maximum length")
	ErrInvalidFormat     = errors.New("invalid format")
	ErrMissingRequired   = errors.New("missing required field")
	ErrMutuallyExclusive = errors.New("mutually exclusive fields set")
	ErrExceedsMaxItems   = errors.New("exceeds maximum items")
	ErrMinItems          = errors.New("requires minimum items")
)

// ValidationError represents a validation failure with context.
type ValidationError struct {
	Field   string
	Value   any
	Message string
	Cause   error
}

func (e ValidationError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Field, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func (e ValidationError) Unwrap() error {
	return e.Cause
}

// newValidationError creates a new ValidationError.
func newValidationError(field, message string, cause error) error {
	return ValidationError{
		Field:   field,
		Message: message,
		Cause:   cause,
	}
}

// validateRequired checks if a string value is non-empty.
func validateRequired(field, value string) error {
	if value == "" {
		return newValidationError(field, "is required", ErrMissingRequired)
	}
	return nil
}

// validateMaxLen checks if a string value exceeds max length.
func validateMaxLen(field, value string, max int) error {
	if len(value) > max {
		return newValidationError(field, fmt.Sprintf("exceeds maximum length of %d characters (got %d)", max, len(value)), ErrExceedsMaxLen)
	}
	return nil
}

// validateRequiredMaxLen validates both required and max length.
func validateRequiredMaxLen(field, value string, max int) error {
	if err := validateRequired(field, value); err != nil {
		return err
	}
	return validateMaxLen(field, value, max)
}

// validateMaxItems checks if a slice exceeds max items.
func validateMaxItems[T any](field string, items []T, max int) error {
	if len(items) > max {
		return newValidationError(field, fmt.Sprintf("exceeds maximum of %d items (got %d)", max, len(items)), ErrExceedsMaxItems)
	}
	return nil
}

// validateMinItems checks if a slice has minimum items.
func validateMinItems[T any](field string, items []T, min int) error {
	if len(items) < min {
		return newValidationError(field, fmt.Sprintf("requires minimum of %d items (got %d)", min, len(items)), ErrMinItems)
	}
	return nil
}
