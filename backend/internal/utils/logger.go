package utils

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

// Sensitive field patterns to redact from logs
var sensitivePatterns = []struct {
	pattern *regexp.Regexp
	label   string
}{
	{regexp.MustCompile(`(?i)"(kata_?sandi[^"]*)":\s*"([^"]*)"`), "kataSandi"},
	{regexp.MustCompile(`(?i)"(password[^"]*)":\s*"([^"]*)"`), "password"},
	{regexp.MustCompile(`(?i)"(pin[^"]*)":\s*"([^"]*)"`), "pin"},
	{regexp.MustCompile(`(?i)"(token[^"]*)":\s*"([^"]*)"`), "token"},
	{regexp.MustCompile(`(?i)"(secret[^"]*)":\s*"([^"]*)"`), "secret"},
	{regexp.MustCompile(`(?i)"(nik)":\s*"([^"]*)"`), "nik"},
	{regexp.MustCompile(`(?i)"(jwt)":\s*"([^"]*)"`), "jwt"},
	{regexp.MustCompile(`(?i)"(authorization)":\s*"([^"]*)"`), "authorization"},
}

// Logger is a simple logger wrapper for services with PII redaction
type Logger struct {
	serviceName string
	logger      *log.Logger
}

// NewLogger creates a new logger instance for a service
func NewLogger(serviceName string) *Logger {
	return &Logger{
		serviceName: serviceName,
		logger:      log.New(os.Stdout, "", log.LstdFlags),
	}
}

// redactSensitiveData removes or masks sensitive information from log data
func redactSensitiveData(data map[string]interface{}) map[string]interface{} {
	if data == nil {
		return nil
	}

	redacted := make(map[string]interface{})
	for key, value := range data {
		lowerKey := strings.ToLower(key)

		// Check if key contains sensitive field names
		if strings.Contains(lowerKey, "password") ||
			strings.Contains(lowerKey, "katasandi") ||
			strings.Contains(lowerKey, "pin") ||
			strings.Contains(lowerKey, "token") ||
			strings.Contains(lowerKey, "secret") ||
			strings.Contains(lowerKey, "jwt") ||
			lowerKey == "nik" {
			redacted[key] = "[REDACTED]"
		} else {
			// For nested maps, recursively redact
			if nestedMap, ok := value.(map[string]interface{}); ok {
				redacted[key] = redactSensitiveData(nestedMap)
			} else {
				redacted[key] = value
			}
		}
	}
	return redacted
}

// formatMessage formats the log message with redacted data
func formatMessage(data map[string]interface{}) string {
	if data == nil {
		return ""
	}

	// Convert to string and apply pattern-based redaction
	msg := fmt.Sprintf("%+v", redactSensitiveData(data))

	// Apply regex patterns for additional protection
	for _, sp := range sensitivePatterns {
		msg = sp.pattern.ReplaceAllString(msg, `"$1":"[REDACTED]"`)
	}

	return msg
}

// Debug logs a debug message with PII redaction
func (l *Logger) Debug(method, message string, data map[string]interface{}) {
	l.logger.Printf("[DEBUG] [%s.%s] %s %s", l.serviceName, method, message, formatMessage(data))
}

// Info logs an info message with PII redaction
func (l *Logger) Info(method, message string, data map[string]interface{}) {
	l.logger.Printf("[INFO] [%s.%s] %s %s", l.serviceName, method, message, formatMessage(data))
}

// Error logs an error message with PII redaction
func (l *Logger) Error(method, message string, err error, data map[string]interface{}) {
	if err != nil {
		l.logger.Printf("[ERROR] [%s.%s] %s: %v %s", l.serviceName, method, message, err, formatMessage(data))
	} else {
		l.logger.Printf("[ERROR] [%s.%s] %s %s", l.serviceName, method, message, formatMessage(data))
	}
}
