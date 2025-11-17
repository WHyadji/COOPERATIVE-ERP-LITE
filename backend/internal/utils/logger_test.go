package utils

import (
	"os"
	"strings"
	"testing"
)

// TestRedactPII tests the PII redaction functionality
func TestRedactPII(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "NIK redaction",
			input:    "User NIK: 1234567890123456",
			expected: "User NIK: [REDACTED-NIK]",
		},
		{
			name:     "Multiple NIKs",
			input:    "NIK 1234567890123456 and 9876543210987654",
			expected: "NIK [REDACTED-NIK] and [REDACTED-NIK]",
		},
		{
			name:     "Password redaction with colon",
			input:    "password: mysecretpass123",
			expected: "password: [REDACTED]",
		},
		{
			name:     "Password redaction with equals",
			input:    "password=supersecret",
			expected: "password: [REDACTED]",
		},
		{
			name:     "Token redaction",
			input:    "token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			expected: "token: [REDACTED]",
		},
		{
			name:     "Secret redaction",
			input:    "secret=my_secret_key_123",
			expected: "secret: [REDACTED]",
		},
		{
			name:     "Email redaction",
			input:    "User email: john.doe@example.com",
			expected: "User email: [REDACTED-EMAIL]",
		},
		{
			name:     "Multiple emails",
			input:    "From: alice@test.com To: bob@example.org",
			expected: "From: [REDACTED-EMAIL] To: [REDACTED-EMAIL]",
		},
		{
			name:     "Indonesian phone number +62",
			input:    "Phone: +6281234567890",
			expected: "Phone: +[REDACTED-PHONE]",
		},
		{
			name:     "Indonesian phone number 62",
			input:    "Contact: 628123456789",
			expected: "Contact: [REDACTED-PHONE]",
		},
		{
			name:     "Indonesian phone number 0",
			input:    "Mobile: 081234567890",
			expected: "Mobile: [REDACTED-PHONE]",
		},
		{
			name:     "Bank account number",
			input:    "Account: 1234567890123",
			expected: "Account: [REDACTED-ACCOUNT]",
		},
		{
			name:     "Short number not redacted",
			input:    "ID: 12345",
			expected: "ID: 12345",
		},
		{
			name:     "Mixed PII",
			input:    "NIK 1234567890123456, email test@example.com, password: secret123, phone +628123456789",
			expected: "NIK [REDACTED-NIK], email [REDACTED-EMAIL], password: [REDACTED], phone +[REDACTED-PHONE]",
		},
		{
			name:     "No PII",
			input:    "This is a normal log message without any sensitive data",
			expected: "This is a normal log message without any sensitive data",
		},
		{
			name:     "JSON with password",
			input:    `{"username":"john","password":"secret123"}`,
			expected: `{"username":"john","password": [REDACTED]"}`,
		},
		{
			name:     "Case insensitive password",
			input:    "Password: MyPassword123",
			expected: "Password: [REDACTED]",
		},
		{
			name:     "PWD variant",
			input:    "pwd=testpass",
			expected: "pwd: [REDACTED]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RedactPII(tt.input)
			if result != tt.expected {
				t.Errorf("RedactPII() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestLogLevel tests log level string conversion
func TestLogLevel(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{LogLevelDebug, "DEBUG"},
		{LogLevelInfo, "INFO"},
		{LogLevelWarn, "WARN"},
		{LogLevelError, "ERROR"},
		{LogLevel(999), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.level.String()
			if result != tt.expected {
				t.Errorf("LogLevel.String() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// TestSetLogLevel tests setting the log level
func TestSetLogLevel(t *testing.T) {
	originalLevel := GetLogLevel()
	defer SetLogLevel(originalLevel) // Restore original level

	SetLogLevel(LogLevelDebug)
	if GetLogLevel() != LogLevelDebug {
		t.Errorf("GetLogLevel() = %v, expected %v", GetLogLevel(), LogLevelDebug)
	}

	SetLogLevel(LogLevelError)
	if GetLogLevel() != LogLevelError {
		t.Errorf("GetLogLevel() = %v, expected %v", GetLogLevel(), LogLevelError)
	}
}

// TestInitLoggerFromEnv tests logger initialization from environment
func TestInitLoggerFromEnv(t *testing.T) {
	originalLevel := GetLogLevel()
	defer SetLogLevel(originalLevel) // Restore original level

	tests := []struct {
		name     string
		envValue string
		expected LogLevel
	}{
		{"debug", "debug", LogLevelDebug},
		{"DEBUG uppercase", "DEBUG", LogLevelDebug},
		{"info", "info", LogLevelInfo},
		{"INFO uppercase", "INFO", LogLevelInfo},
		{"warn", "warn", LogLevelWarn},
		{"WARN uppercase", "WARN", LogLevelWarn},
		{"warning", "warning", LogLevelWarn},
		{"error", "error", LogLevelError},
		{"ERROR uppercase", "ERROR", LogLevelError},
		{"invalid", "invalid", LogLevelInfo}, // Default to INFO
		{"empty", "", LogLevelInfo},          // Default to INFO
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = os.Setenv("LOG_LEVEL", tt.envValue)
			defer func() { _ = os.Unsetenv("LOG_LEVEL") }()

			InitLoggerFromEnv()

			if GetLogLevel() != tt.expected {
				t.Errorf("InitLoggerFromEnv() set level to %v, expected %v", GetLogLevel(), tt.expected)
			}
		})
	}
}

// TestShouldLog tests the log filtering logic
func TestShouldLog(t *testing.T) {
	originalLevel := GetLogLevel()
	defer SetLogLevel(originalLevel)

	tests := []struct {
		name          string
		currentLevel  LogLevel
		messageLevel  LogLevel
		shouldDisplay bool
	}{
		{"Debug level shows debug", LogLevelDebug, LogLevelDebug, true},
		{"Debug level shows info", LogLevelDebug, LogLevelInfo, true},
		{"Debug level shows warn", LogLevelDebug, LogLevelWarn, true},
		{"Debug level shows error", LogLevelDebug, LogLevelError, true},
		{"Info level hides debug", LogLevelInfo, LogLevelDebug, false},
		{"Info level shows info", LogLevelInfo, LogLevelInfo, true},
		{"Info level shows warn", LogLevelInfo, LogLevelWarn, true},
		{"Info level shows error", LogLevelInfo, LogLevelError, true},
		{"Warn level hides debug", LogLevelWarn, LogLevelDebug, false},
		{"Warn level hides info", LogLevelWarn, LogLevelInfo, false},
		{"Warn level shows warn", LogLevelWarn, LogLevelWarn, true},
		{"Warn level shows error", LogLevelWarn, LogLevelError, true},
		{"Error level hides debug", LogLevelError, LogLevelDebug, false},
		{"Error level hides info", LogLevelError, LogLevelInfo, false},
		{"Error level hides warn", LogLevelError, LogLevelWarn, false},
		{"Error level shows error", LogLevelError, LogLevelError, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLogLevel(tt.currentLevel)
			result := shouldLog(tt.messageLevel)

			if result != tt.shouldDisplay {
				t.Errorf("shouldLog() = %v, expected %v (current=%v, message=%v)",
					result, tt.shouldDisplay, tt.currentLevel, tt.messageLevel)
			}
		})
	}
}

// TestLogFunctions tests that log functions don't panic
func TestLogFunctions(t *testing.T) {
	originalLevel := GetLogLevel()
	defer SetLogLevel(originalLevel)

	// Set to DEBUG to ensure all messages are logged
	SetLogLevel(LogLevelDebug)

	tests := []struct {
		name string
		fn   func()
	}{
		{
			name: "LogDebug",
			fn:   func() { LogDebug("debug message") },
		},
		{
			name: "LogInfo",
			fn:   func() { LogInfo("info message") },
		},
		{
			name: "LogWarn",
			fn:   func() { LogWarn("warn message") },
		},
		{
			name: "LogError",
			fn:   func() { LogError("error message", nil) },
		},
		{
			name: "LogErrorf",
			fn:   func() { LogErrorf("formatted error: %s", "test") },
		},
		{
			name: "LogInfof",
			fn:   func() { LogInfof("formatted info: %d", 42) },
		},
		{
			name: "LogDebugf",
			fn:   func() { LogDebugf("formatted debug: %v", true) },
		},
		{
			name: "LogWarnf",
			fn:   func() { LogWarnf("formatted warn: %.2f", 3.14) },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("%s panicked: %v", tt.name, r)
				}
			}()

			tt.fn()
		})
	}
}

// TestRedactPIIEdgeCases tests edge cases in PII redaction
func TestRedactPIIEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Only spaces",
			input:    "     ",
			expected: "     ",
		},
		{
			name:     "NIK with prefix",
			input:    "NIK:1234567890123456",
			expected: "NIK:[REDACTED-NIK]",
		},
		{
			name:     "Phone with parentheses",
			input:    "Call (+62)81234567890 now",
			expected: "Call (+62)[REDACTED-ACCOUNT] now",
		},
		{
			name:     "Email in quotes",
			input:    `"email":"test@example.com"`,
			expected: `"email":"[REDACTED-EMAIL]"`,
		},
		{
			name:     "Password with special chars",
			input:    "password: P@ssw0rd!#$",
			expected: "password: [REDACTED]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RedactPII(tt.input)
			if result != tt.expected {
				t.Errorf("RedactPII() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

// BenchmarkRedactPII benchmarks the PII redaction function
func BenchmarkRedactPII(b *testing.B) {
	testMessage := "User with NIK 1234567890123456, email john@example.com, password: secret123, phone +628123456789 logged in"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RedactPII(testMessage)
	}
}

// BenchmarkRedactPIINoMatch benchmarks redaction with no PII
func BenchmarkRedactPIINoMatch(b *testing.B) {
	testMessage := "This is a normal log message without any sensitive data at all"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		RedactPII(testMessage)
	}
}

// BenchmarkLogInfo benchmarks info logging
func BenchmarkLogInfo(b *testing.B) {
	originalLevel := GetLogLevel()
	defer SetLogLevel(originalLevel)
	SetLogLevel(LogLevelInfo)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LogInfo("test message")
	}
}

// TestFormatLogMessage tests message formatting
func TestFormatLogMessage(t *testing.T) {
	msg := formatLogMessage(LogLevelInfo, "test message", nil)

	// Check format: [timestamp] LEVEL: message
	if !strings.Contains(msg, "INFO") {
		t.Error("Formatted message should contain log level")
	}

	if !strings.Contains(msg, "test message") {
		t.Error("Formatted message should contain the message")
	}

	// Check timestamp format (should start with [YYYY-MM-DD)
	if !strings.HasPrefix(msg, "[202") && !strings.HasPrefix(msg, "[203") {
		t.Error("Formatted message should start with timestamp")
	}
}
