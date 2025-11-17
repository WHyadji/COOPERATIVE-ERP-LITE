package utils

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

// WrapDatabaseError wraps database errors with user-friendly messages
// Uses fmt.Errorf with %w to preserve error chain for errors.Is() checks
func WrapDatabaseError(err error, context string) error {
	if err == nil {
		return nil
	}

	// Handle specific GORM errors
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("%s tidak ditemukan: %w", context, err)
	}

	if errors.Is(err, gorm.ErrInvalidData) {
		return fmt.Errorf("data tidak valid: %w", err)
	}

	if errors.Is(err, gorm.ErrInvalidField) {
		return fmt.Errorf("field tidak valid: %w", err)
	}

	// Default error message with wrapped error
	return fmt.Errorf("%s: %w", context, err)
}

// NewValidationError creates a validation error
func NewValidationError(message string) error {
	return errors.New(message)
}

// Common error types
var (
	ErrInvalidAmount = errors.New("jumlah tidak valid")
)

// WrapValidationError wraps a validation error with additional context
// Note: second parameter intentionally unused (reserved for future use)
func WrapValidationError(err error, _ interface{}) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("validation error: %w", err)
}

// WrapGenerationError wraps an error that occurred during generation
func WrapGenerationError(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}
