package middleware

import (
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiterType represents the type of rate limiting algorithm
type RateLimiterType string

const (
	TokenBucketType   RateLimiterType = "token_bucket"
	SlidingWindowType RateLimiterType = "sliding_window"
)

// RateLimiterConfig holds configuration for rate limiters
type RateLimiterConfig struct {
	Algorithm         RateLimiterType
	RequestsPerMinute int
	BurstSize         int  // Only for token bucket
	Window            time.Duration // Only for sliding window
}

// DefaultRateLimiterConfig returns the default configuration
// Uses Token Bucket algorithm by default for better performance
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		Algorithm:         TokenBucketType,
		RequestsPerMinute: 100,
		BurstSize:         20,
		Window:            1 * time.Minute,
	}
}

// LoadRateLimiterConfig loads configuration from environment variables
// Environment variables:
// - RATE_LIMIT_ALGORITHM: "token_bucket" or "sliding_window" (default: token_bucket)
// - RATE_LIMIT_RPM: Requests per minute (default: 100)
// - RATE_LIMIT_BURST: Burst size for token bucket (default: 20)
func LoadRateLimiterConfig() RateLimiterConfig {
	cfg := DefaultRateLimiterConfig()

	// Load algorithm type
	if algo := os.Getenv("RATE_LIMIT_ALGORITHM"); algo != "" {
		cfg.Algorithm = RateLimiterType(algo)
	}

	// Load requests per minute
	if rpm := os.Getenv("RATE_LIMIT_RPM"); rpm != "" {
		if val, err := strconv.Atoi(rpm); err == nil {
			cfg.RequestsPerMinute = val
		}
	}

	// Load burst size
	if burst := os.Getenv("RATE_LIMIT_BURST"); burst != "" {
		if val, err := strconv.Atoi(burst); err == nil {
			cfg.BurstSize = val
		}
	}

	return cfg
}

// NewRateLimiterMiddleware creates a rate limiter middleware based on configuration
// Supports both Token Bucket and Sliding Window algorithms
func NewRateLimiterMiddleware(cfg RateLimiterConfig) gin.HandlerFunc {
	switch cfg.Algorithm {
	case TokenBucketType:
		return TokenBucketMiddleware(cfg.RequestsPerMinute, cfg.BurstSize)
	case SlidingWindowType:
		return RateLimitMiddleware(cfg.RequestsPerMinute, cfg.Window)
	default:
		// Default to token bucket if unknown algorithm
		return TokenBucketMiddleware(cfg.RequestsPerMinute, cfg.BurstSize)
	}
}

// NewRateLimiterMiddlewareFromEnv creates a rate limiter middleware from environment variables
// This is the recommended way to create rate limiters in production
func NewRateLimiterMiddlewareFromEnv() gin.HandlerFunc {
	cfg := LoadRateLimiterConfig()
	return NewRateLimiterMiddleware(cfg)
}

// LoginRateLimiterConfig holds configuration for login rate limiters
type LoginRateLimiterConfig struct {
	Algorithm       RateLimiterType
	MaxAttempts     int
	Window          time.Duration
	LockoutDuration time.Duration
}

// DefaultLoginRateLimiterConfig returns the default configuration for login endpoints
func DefaultLoginRateLimiterConfig() LoginRateLimiterConfig {
	return LoginRateLimiterConfig{
		Algorithm:       TokenBucketType,
		MaxAttempts:     5,
		Window:          15 * time.Minute,
		LockoutDuration: 15 * time.Minute,
	}
}

// LoadLoginRateLimiterConfig loads login rate limiter config from environment
// Environment variables:
// - LOGIN_RATE_LIMIT_ALGORITHM: "token_bucket" or "sliding_window"
// - LOGIN_MAX_ATTEMPTS: Maximum failed attempts (default: 5)
// - LOGIN_WINDOW_MINUTES: Time window in minutes (default: 15)
// - LOGIN_LOCKOUT_MINUTES: Lockout duration in minutes (default: 15)
func LoadLoginRateLimiterConfig() LoginRateLimiterConfig {
	cfg := DefaultLoginRateLimiterConfig()

	// Load algorithm
	if algo := os.Getenv("LOGIN_RATE_LIMIT_ALGORITHM"); algo != "" {
		cfg.Algorithm = RateLimiterType(algo)
	}

	// Load max attempts
	if attempts := os.Getenv("LOGIN_MAX_ATTEMPTS"); attempts != "" {
		if val, err := strconv.Atoi(attempts); err == nil {
			cfg.MaxAttempts = val
		}
	}

	// Load window
	if window := os.Getenv("LOGIN_WINDOW_MINUTES"); window != "" {
		if val, err := strconv.Atoi(window); err == nil {
			cfg.Window = time.Duration(val) * time.Minute
		}
	}

	// Load lockout duration
	if lockout := os.Getenv("LOGIN_LOCKOUT_MINUTES"); lockout != "" {
		if val, err := strconv.Atoi(lockout); err == nil {
			cfg.LockoutDuration = time.Duration(val) * time.Minute
		}
	}

	return cfg
}

// NewLoginRateLimiterMiddleware creates a login rate limiter middleware based on configuration
func NewLoginRateLimiterMiddleware(cfg LoginRateLimiterConfig) gin.HandlerFunc {
	switch cfg.Algorithm {
	case TokenBucketType:
		return LoginTokenBucketMiddleware(cfg.MaxAttempts, cfg.Window, cfg.LockoutDuration)
	case SlidingWindowType:
		return LoginRateLimitMiddleware(cfg.MaxAttempts, cfg.Window, cfg.LockoutDuration)
	default:
		// Default to token bucket
		return LoginTokenBucketMiddleware(cfg.MaxAttempts, cfg.Window, cfg.LockoutDuration)
	}
}

// NewLoginRateLimiterMiddlewareFromEnv creates a login rate limiter from environment variables
func NewLoginRateLimiterMiddlewareFromEnv() gin.HandlerFunc {
	cfg := LoadLoginRateLimiterConfig()
	return NewLoginRateLimiterMiddleware(cfg)
}

// Example usage in main.go:
//
// func setupRouter() *gin.Engine {
//     router := gin.Default()
//
//     // General API rate limiting (from environment)
//     router.Use(middleware.NewRateLimiterMiddlewareFromEnv())
//
//     // Or with custom config:
//     // cfg := middleware.RateLimiterConfig{
//     //     Algorithm: middleware.TokenBucketType,
//     //     RequestsPerMinute: 100,
//     //     BurstSize: 20,
//     // }
//     // router.Use(middleware.NewRateLimiterMiddleware(cfg))
//
//     // Login endpoint (from environment)
//     authGroup := router.Group("/auth")
//     authGroup.Use(middleware.NewLoginRateLimiterMiddlewareFromEnv())
//
//     return router
// }
//
// Environment configuration example:
//
// # .env or fly.toml secrets
// RATE_LIMIT_ALGORITHM=token_bucket
// RATE_LIMIT_RPM=100
// RATE_LIMIT_BURST=20
//
// LOGIN_RATE_LIMIT_ALGORITHM=token_bucket
// LOGIN_MAX_ATTEMPTS=5
// LOGIN_WINDOW_MINUTES=15
// LOGIN_LOCKOUT_MINUTES=15
//
// # To switch algorithms without code change:
// RATE_LIMIT_ALGORITHM=sliding_window  # Uses old algorithm
// RATE_LIMIT_ALGORITHM=token_bucket    # Uses new algorithm
