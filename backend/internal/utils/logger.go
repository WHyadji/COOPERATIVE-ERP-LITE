package utils

import (
	"log"
	"os"
)

// Logger is a simple logger wrapper for services
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

// Debug logs a debug message
func (l *Logger) Debug(method, message string, data map[string]interface{}) {
	l.logger.Printf("[DEBUG] [%s.%s] %s %+v", l.serviceName, method, message, data)
}

// Info logs an info message
func (l *Logger) Info(method, message string, data map[string]interface{}) {
	l.logger.Printf("[INFO] [%s.%s] %s %+v", l.serviceName, method, message, data)
}

// Error logs an error message
func (l *Logger) Error(method, message string, err error, data map[string]interface{}) {
	if err != nil {
		l.logger.Printf("[ERROR] [%s.%s] %s: %v %+v", l.serviceName, method, message, err, data)
	} else {
		l.logger.Printf("[ERROR] [%s.%s] %s %+v", l.serviceName, method, message, data)
	}
}
