package utils

import (
	"errors"
	"gorm.io/gorm"
)

// WrapDatabaseError wraps database errors with user-friendly messages
func WrapDatabaseError(err error, context string) error {
	if err == nil {
		return nil
	}

	// Handle specific GORM errors
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(context + " tidak ditemukan")
	}

	if errors.Is(err, gorm.ErrInvalidData) {
		return errors.New("data tidak valid: " + err.Error())
	}

	if errors.Is(err, gorm.ErrInvalidField) {
		return errors.New("field tidak valid: " + err.Error())
	}

	// Default error message
	return errors.New(context + ": " + err.Error())
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
func WrapValidationError(err error, context interface{}) error {
	if err == nil {
		return nil
	}
	return errors.New(err.Error())
}

// WrapGenerationError wraps an error that occurred during generation
func WrapGenerationError(err error, context string) error {
	if err == nil {
		return nil
	}
	return errors.New(context + ": " + err.Error())
}
