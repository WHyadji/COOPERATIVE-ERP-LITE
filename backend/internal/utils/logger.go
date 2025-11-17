package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger configuration
var (
	currentLogLevel = LogLevelInfo // Default log level

	// PII patterns for redaction (extended set)
	nikPattern      = regexp.MustCompile(`\b\d{16}\b`)                                    // 16-digit NIK
	passwordPattern = regexp.MustCompile(`(?i)(password|passwd|pwd|secret|token)["']?\s*[:=]\s*["']?[^"'\s,}]+`) // Passwords, tokens
	emailPattern    = regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b`) // Email addresses
	phonePattern    = regexp.MustCompile(`\b(\+62|62|0)[0-9]{8,12}\b`)                   // Indonesian phone numbers
	bankAcctPattern = regexp.MustCompile(`\b\d{10,16}\b`)                                // Bank account numbers (10-16 digits)
)

// SetLogLevel sets the minimum log level that will be output
// Only logs with this level or higher will be printed
func SetLogLevel(level LogLevel) {
	currentLogLevel = level
}

// GetLogLevel returns the current log level
func GetLogLevel() LogLevel {
	return currentLogLevel
}

// InitLoggerFromEnv initializes the logger based on environment variables
// Reads LOG_LEVEL environment variable: "debug", "info", "warn", "error"
func InitLoggerFromEnv() {
	logLevelStr := strings.ToLower(os.Getenv("LOG_LEVEL"))

	switch logLevelStr {
	case "debug":
		SetLogLevel(LogLevelDebug)
	case "info":
		SetLogLevel(LogLevelInfo)
	case "warn", "warning":
		SetLogLevel(LogLevelWarn)
	case "error":
		SetLogLevel(LogLevelError)
	default:
		// Default to INFO if not specified or invalid
		SetLogLevel(LogLevelInfo)
	}
}

// RedactPII removes personally identifiable information from log messages
// This prevents sensitive data from appearing in logs
func RedactPII(message string) string {
	// Redact NIK (Indonesian ID numbers)
	message = nikPattern.ReplaceAllString(message, "[REDACTED-NIK]")

	// Redact passwords and tokens
	message = passwordPattern.ReplaceAllStringFunc(message, func(match string) string {
		parts := regexp.MustCompile(`[:=]`).Split(match, 2)
		if len(parts) == 2 {
			return parts[0] + ": [REDACTED]"
		}
		return "[REDACTED]"
	})

	// Redact email addresses
	message = emailPattern.ReplaceAllString(message, "[REDACTED-EMAIL]")

	// Redact phone numbers
	message = phonePattern.ReplaceAllString(message, "[REDACTED-PHONE]")

	// Redact bank account numbers (but preserve smaller numbers that might be IDs)
	message = bankAcctPattern.ReplaceAllStringFunc(message, func(match string) string {
		// Only redact if it looks like a bank account (10+ digits)
		if len(match) >= 10 {
			return "[REDACTED-ACCOUNT]"
		}
		return match
	})

	return message
}

// formatLogMessage formats a log message with timestamp and level
func formatLogMessage(level LogLevel, message string, err error) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	if err != nil {
		message = fmt.Sprintf("%s - %v", message, err)
	}

	// Redact PII before logging
	message = RedactPII(message)

	return fmt.Sprintf("[%s] %s: %s", timestamp, level.String(), message)
}

// shouldLog checks if a message with the given level should be logged
func shouldLog(level LogLevel) bool {
	return level >= currentLogLevel
}

// LogDebug logs a debug-level message
// Only logged when log level is DEBUG
func LogDebug(message string) {
	if shouldLog(LogLevelDebug) {
		fmt.Println(formatLogMessage(LogLevelDebug, message, nil))
	}
}

// LogInfo logs an info-level message
// Logged when log level is DEBUG or INFO
func LogInfo(message string) {
	if shouldLog(LogLevelInfo) {
		fmt.Println(formatLogMessage(LogLevelInfo, message, nil))
	}
}

// LogWarn logs a warning-level message
// Logged when log level is DEBUG, INFO, or WARN
func LogWarn(message string) {
	if shouldLog(LogLevelWarn) {
		fmt.Println(formatLogMessage(LogLevelWarn, message, nil))
	}
}

// LogError logs an error-level message with optional error object
// Always logged (unless log level is explicitly disabled)
func LogError(message string, err error) {
	if shouldLog(LogLevelError) {
		fmt.Println(formatLogMessage(LogLevelError, message, err))
	}
}

// LogErrorf logs a formatted error message
func LogErrorf(format string, args ...interface{}) {
	if shouldLog(LogLevelError) {
		message := fmt.Sprintf(format, args...)
		fmt.Println(formatLogMessage(LogLevelError, message, nil))
	}
}

// LogInfof logs a formatted info message
func LogInfof(format string, args ...interface{}) {
	if shouldLog(LogLevelInfo) {
		message := fmt.Sprintf(format, args...)
		fmt.Println(formatLogMessage(LogLevelInfo, message, nil))
	}
}

// LogDebugf logs a formatted debug message
func LogDebugf(format string, args ...interface{}) {
	if shouldLog(LogLevelDebug) {
		message := fmt.Sprintf(format, args...)
		fmt.Println(formatLogMessage(LogLevelDebug, message, nil))
	}
}

// LogWarnf logs a formatted warning message
func LogWarnf(format string, args ...interface{}) {
	if shouldLog(LogLevelWarn) {
		message := fmt.Sprintf(format, args...)
		fmt.Println(formatLogMessage(LogLevelWarn, message, nil))
	}
}
