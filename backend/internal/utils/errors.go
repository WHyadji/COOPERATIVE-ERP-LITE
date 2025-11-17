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

// Common error messages (safe for client exposure)
// These messages are designed to be user-friendly and prevent information disclosure
var (
	// Generic messages
	ErrMsgInternalServer  = "Terjadi kesalahan pada server"
	ErrMsgNotFound        = "Data tidak ditemukan"
	ErrMsgInvalidInput    = "Data yang dimasukkan tidak valid"
	ErrMsgUnauthorized    = "Anda tidak memiliki akses"
	ErrMsgDatabaseError   = "Terjadi kesalahan pada database"
	ErrMsgValidationError = "Validasi data gagal"

	// Specific messages
	ErrMsgDuplicateEntry  = "Data sudah ada dalam sistem"
	ErrMsgForeignKeyError = "Data masih digunakan oleh data lain"
	ErrMsgConnectionError = "Tidak dapat terhubung ke database"
	ErrMsgInvalidAmount   = "Jumlah tidak valid"
)

// Common error types (for internal use)
var (
	ErrInvalidAmount = errors.New("jumlah tidak valid")
)

// SanitizeError converts internal errors into safe, user-facing error messages
// This prevents information disclosure by mapping internal errors to generic messages
// Always use this function before sending errors to API clients
func SanitizeError(err error) string {
	if err == nil {
		return ""
	}

	// Check for specific GORM errors
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrMsgNotFound
	}

	if errors.Is(err, gorm.ErrInvalidData) || errors.Is(err, gorm.ErrInvalidField) {
		return ErrMsgInvalidInput
	}

	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrMsgDuplicateEntry
	}

	if errors.Is(err, gorm.ErrForeignKeyViolated) {
		return ErrMsgForeignKeyError
	}

	// Check for internal validation errors
	if errors.Is(err, ErrInvalidAmount) {
		return ErrMsgInvalidAmount
	}

	// Check error message for common database issues (connection, timeout)
	errMsg := err.Error()
	if containsAny(errMsg, []string{"connection", "timeout", "dial", "network"}) {
		return ErrMsgConnectionError
	}

	// Check for validation keywords
	if containsAny(errMsg, []string{"validation", "invalid", "required", "must"}) {
		return ErrMsgValidationError
	}

	// Default to generic internal server error (never expose raw error details)
	return ErrMsgInternalServer
}

// containsAny checks if a string contains any of the provided substrings (case-insensitive)
func containsAny(s string, substrs []string) bool {
	sLower := toLower(s)
	for _, substr := range substrs {
		substrLower := toLower(substr)
		if contains(sLower, substrLower) {
			return true
		}
	}
	return false
}

// toLower converts a string to lowercase without external dependencies
func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result[i] = c + ('a' - 'A')
		} else {
			result[i] = c
		}
	}
	return string(result)
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

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
