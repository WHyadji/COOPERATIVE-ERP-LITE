package utils

import (
	"errors"
	"testing"

	"gorm.io/gorm"
)

// TestSanitizeError tests the error sanitization function
func TestSanitizeError(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected string
	}{
		{
			name:     "nil error",
			input:    nil,
			expected: "",
		},
		{
			name:     "GORM record not found",
			input:    gorm.ErrRecordNotFound,
			expected: ErrMsgNotFound,
		},
		{
			name:     "GORM invalid data",
			input:    gorm.ErrInvalidData,
			expected: ErrMsgInvalidInput,
		},
		{
			name:     "GORM invalid field",
			input:    gorm.ErrInvalidField,
			expected: ErrMsgInvalidInput,
		},
		{
			name:     "GORM duplicated key",
			input:    gorm.ErrDuplicatedKey,
			expected: ErrMsgDuplicateEntry,
		},
		{
			name:     "GORM foreign key violated",
			input:    gorm.ErrForeignKeyViolated,
			expected: ErrMsgForeignKeyError,
		},
		{
			name:     "Internal validation error",
			input:    ErrInvalidAmount,
			expected: ErrMsgInvalidAmount,
		},
		{
			name:     "Connection error",
			input:    errors.New("connection refused"),
			expected: ErrMsgConnectionError,
		},
		{
			name:     "Timeout error",
			input:    errors.New("context deadline exceeded (timeout)"),
			expected: ErrMsgConnectionError,
		},
		{
			name:     "Network error",
			input:    errors.New("network unreachable"),
			expected: ErrMsgConnectionError,
		},
		{
			name:     "Dial error",
			input:    errors.New("dial tcp: connection refused"),
			expected: ErrMsgConnectionError,
		},
		{
			name:     "Validation error",
			input:    errors.New("validation failed: field is required"),
			expected: ErrMsgValidationError,
		},
		{
			name:     "Invalid input error",
			input:    errors.New("invalid email format"),
			expected: ErrMsgValidationError,
		},
		{
			name:     "Required field error",
			input:    errors.New("name is required"),
			expected: ErrMsgValidationError,
		},
		{
			name:     "Must constraint error",
			input:    errors.New("password must be at least 8 characters"),
			expected: ErrMsgValidationError,
		},
		{
			name:     "Generic internal error",
			input:    errors.New("unexpected nil pointer dereference"),
			expected: ErrMsgInternalServer,
		},
		{
			name:     "Database constraint error",
			input:    errors.New("pq: duplicate key value violates unique constraint"),
			expected: ErrMsgInternalServer,
		},
		{
			name:     "File system error",
			input:    errors.New("no such file or directory: /path/to/secret"),
			expected: ErrMsgInternalServer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeError(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeError() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestWrapDatabaseError tests the database error wrapping function
func TestWrapDatabaseError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		context  string
		expected string
		isNil    bool
	}{
		{
			name:     "nil error",
			err:      nil,
			context:  "test context",
			expected: "",
			isNil:    true,
		},
		{
			name:     "record not found",
			err:      gorm.ErrRecordNotFound,
			context:  "Anggota",
			expected: "Anggota tidak ditemukan",
			isNil:    false,
		},
		{
			name:     "invalid data",
			err:      gorm.ErrInvalidData,
			context:  "Simpanan",
			expected: "data tidak valid",
			isNil:    false,
		},
		{
			name:     "invalid field",
			err:      gorm.ErrInvalidField,
			context:  "Transaksi",
			expected: "field tidak valid",
			isNil:    false,
		},
		{
			name:     "generic error",
			err:      errors.New("something went wrong"),
			context:  "Database operation",
			expected: "Database operation: something went wrong",
			isNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapDatabaseError(tt.err, tt.context)

			if tt.isNil && result != nil {
				t.Errorf("WrapDatabaseError() should return nil, got %v", result)
				return
			}

			if !tt.isNil && result == nil {
				t.Errorf("WrapDatabaseError() should not return nil")
				return
			}

			if result != nil && !contains(result.Error(), tt.expected) {
				t.Errorf("WrapDatabaseError() error message should contain %q, got %q", tt.expected, result.Error())
			}

			// Verify error chain is preserved
			if !tt.isNil && tt.err != nil && !errors.Is(result, tt.err) {
				t.Errorf("WrapDatabaseError() should preserve error chain for errors.Is()")
			}
		})
	}
}

// TestWrapValidationError tests the validation error wrapping function
func TestWrapValidationError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		context  interface{}
		expected string
		isNil    bool
	}{
		{
			name:     "nil error",
			err:      nil,
			context:  nil,
			expected: "",
			isNil:    true,
		},
		{
			name:     "validation error with context",
			err:      errors.New("field is required"),
			context:  "user input",
			expected: "validation error: field is required",
			isNil:    false,
		},
		{
			name:     "validation error without context",
			err:      errors.New("invalid format"),
			context:  nil,
			expected: "validation error: invalid format",
			isNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapValidationError(tt.err, tt.context)

			if tt.isNil && result != nil {
				t.Errorf("WrapValidationError() should return nil, got %v", result)
				return
			}

			if !tt.isNil && result == nil {
				t.Errorf("WrapValidationError() should not return nil")
				return
			}

			if result != nil && result.Error() != tt.expected {
				t.Errorf("WrapValidationError() = %q, expected %q", result.Error(), tt.expected)
			}
		})
	}
}

// TestWrapGenerationError tests the generation error wrapping function
func TestWrapGenerationError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		context  string
		expected string
		isNil    bool
	}{
		{
			name:     "nil error",
			err:      nil,
			context:  "test context",
			expected: "",
			isNil:    true,
		},
		{
			name:     "generation error",
			err:      errors.New("failed to generate report"),
			context:  "Laporan Keuangan",
			expected: "Laporan Keuangan: failed to generate report",
			isNil:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapGenerationError(tt.err, tt.context)

			if tt.isNil && result != nil {
				t.Errorf("WrapGenerationError() should return nil, got %v", result)
				return
			}

			if !tt.isNil && result == nil {
				t.Errorf("WrapGenerationError() should not return nil")
				return
			}

			if result != nil && result.Error() != tt.expected {
				t.Errorf("WrapGenerationError() = %q, expected %q", result.Error(), tt.expected)
			}
		})
	}
}

// TestNewValidationError tests the validation error creation function
func TestNewValidationError(t *testing.T) {
	message := "test validation error"
	err := NewValidationError(message)

	if err == nil {
		t.Fatal("NewValidationError() should not return nil")
	}

	if err.Error() != message {
		t.Errorf("NewValidationError() = %q, expected %q", err.Error(), message)
	}
}

// TestErrorConstants tests that error constants are properly defined
func TestErrorConstants(t *testing.T) {
	if ErrInvalidAmount == nil {
		t.Error("ErrInvalidAmount should not be nil")
	}

	if ErrInvalidAmount.Error() != "jumlah tidak valid" {
		t.Errorf("ErrInvalidAmount message = %q, expected \"jumlah tidak valid\"", ErrInvalidAmount.Error())
	}
}

// TestErrorMessageConstants tests that error message constants are properly defined
func TestErrorMessageConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"ErrMsgInternalServer", ErrMsgInternalServer, "Terjadi kesalahan pada server"},
		{"ErrMsgNotFound", ErrMsgNotFound, "Data tidak ditemukan"},
		{"ErrMsgInvalidInput", ErrMsgInvalidInput, "Data yang dimasukkan tidak valid"},
		{"ErrMsgUnauthorized", ErrMsgUnauthorized, "Anda tidak memiliki akses"},
		{"ErrMsgDatabaseError", ErrMsgDatabaseError, "Terjadi kesalahan pada database"},
		{"ErrMsgValidationError", ErrMsgValidationError, "Validasi data gagal"},
		{"ErrMsgDuplicateEntry", ErrMsgDuplicateEntry, "Data sudah ada dalam sistem"},
		{"ErrMsgForeignKeyError", ErrMsgForeignKeyError, "Data masih digunakan oleh data lain"},
		{"ErrMsgConnectionError", ErrMsgConnectionError, "Tidak dapat terhubung ke database"},
		{"ErrMsgInvalidAmount", ErrMsgInvalidAmount, "Jumlah tidak valid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("%s = %q, expected %q", tt.name, tt.constant, tt.expected)
			}
		})
	}
}

// Benchmark tests
func BenchmarkSanitizeError(b *testing.B) {
	testErr := errors.New("test error: connection timeout occurred")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SanitizeError(testErr)
	}
}

func BenchmarkSanitizeErrorGORM(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SanitizeError(gorm.ErrRecordNotFound)
	}
}
