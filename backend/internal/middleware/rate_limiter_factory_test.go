package middleware

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultRateLimiterConfig(t *testing.T) {
	cfg := DefaultRateLimiterConfig()

	assert.Equal(t, TokenBucketType, cfg.Algorithm)
	assert.Equal(t, 100, cfg.RequestsPerMinute)
	assert.Equal(t, 20, cfg.BurstSize)
	assert.Equal(t, 1*time.Minute, cfg.Window)
}

func TestLoadRateLimiterConfig(t *testing.T) {
	t.Run("loads from environment variables", func(t *testing.T) {
		// Set env vars
		os.Setenv("RATE_LIMIT_ALGORITHM", "sliding_window")
		os.Setenv("RATE_LIMIT_RPM", "150")
		os.Setenv("RATE_LIMIT_BURST", "30")
		defer func() {
			os.Unsetenv("RATE_LIMIT_ALGORITHM")
			os.Unsetenv("RATE_LIMIT_RPM")
			os.Unsetenv("RATE_LIMIT_BURST")
		}()

		cfg := LoadRateLimiterConfig()

		assert.Equal(t, SlidingWindowType, cfg.Algorithm)
		assert.Equal(t, 150, cfg.RequestsPerMinute)
		assert.Equal(t, 30, cfg.BurstSize)
	})

	t.Run("uses defaults when env vars not set", func(t *testing.T) {
		cfg := LoadRateLimiterConfig()

		assert.Equal(t, TokenBucketType, cfg.Algorithm)
		assert.Equal(t, 100, cfg.RequestsPerMinute)
		assert.Equal(t, 20, cfg.BurstSize)
	})

	t.Run("handles invalid env values gracefully", func(t *testing.T) {
		os.Setenv("RATE_LIMIT_RPM", "invalid")
		defer os.Unsetenv("RATE_LIMIT_RPM")

		cfg := LoadRateLimiterConfig()

		// Should use default value
		assert.Equal(t, 100, cfg.RequestsPerMinute)
	})
}

func TestNewRateLimiterMiddleware(t *testing.T) {
	t.Run("creates token bucket middleware", func(t *testing.T) {
		cfg := RateLimiterConfig{
			Algorithm:         TokenBucketType,
			RequestsPerMinute: 100,
			BurstSize:         20,
		}

		middleware := NewRateLimiterMiddleware(cfg)

		assert.NotNil(t, middleware)
	})

	t.Run("creates sliding window middleware", func(t *testing.T) {
		cfg := RateLimiterConfig{
			Algorithm:         SlidingWindowType,
			RequestsPerMinute: 100,
			Window:            1 * time.Minute,
		}

		middleware := NewRateLimiterMiddleware(cfg)

		assert.NotNil(t, middleware)
	})

	t.Run("defaults to token bucket for unknown algorithm", func(t *testing.T) {
		cfg := RateLimiterConfig{
			Algorithm:         "unknown",
			RequestsPerMinute: 100,
			BurstSize:         20,
		}

		middleware := NewRateLimiterMiddleware(cfg)

		assert.NotNil(t, middleware)
	})
}

func TestDefaultLoginRateLimiterConfig(t *testing.T) {
	cfg := DefaultLoginRateLimiterConfig()

	assert.Equal(t, TokenBucketType, cfg.Algorithm)
	assert.Equal(t, 5, cfg.MaxAttempts)
	assert.Equal(t, 15*time.Minute, cfg.Window)
	assert.Equal(t, 15*time.Minute, cfg.LockoutDuration)
}

func TestLoadLoginRateLimiterConfig(t *testing.T) {
	t.Run("loads from environment variables", func(t *testing.T) {
		os.Setenv("LOGIN_RATE_LIMIT_ALGORITHM", "sliding_window")
		os.Setenv("LOGIN_MAX_ATTEMPTS", "3")
		os.Setenv("LOGIN_WINDOW_MINUTES", "10")
		os.Setenv("LOGIN_LOCKOUT_MINUTES", "20")
		defer func() {
			os.Unsetenv("LOGIN_RATE_LIMIT_ALGORITHM")
			os.Unsetenv("LOGIN_MAX_ATTEMPTS")
			os.Unsetenv("LOGIN_WINDOW_MINUTES")
			os.Unsetenv("LOGIN_LOCKOUT_MINUTES")
		}()

		cfg := LoadLoginRateLimiterConfig()

		assert.Equal(t, SlidingWindowType, cfg.Algorithm)
		assert.Equal(t, 3, cfg.MaxAttempts)
		assert.Equal(t, 10*time.Minute, cfg.Window)
		assert.Equal(t, 20*time.Minute, cfg.LockoutDuration)
	})

	t.Run("uses defaults when env vars not set", func(t *testing.T) {
		cfg := LoadLoginRateLimiterConfig()

		assert.Equal(t, TokenBucketType, cfg.Algorithm)
		assert.Equal(t, 5, cfg.MaxAttempts)
		assert.Equal(t, 15*time.Minute, cfg.Window)
	})
}

func TestNewLoginRateLimiterMiddleware(t *testing.T) {
	t.Run("creates token bucket login middleware", func(t *testing.T) {
		cfg := LoginRateLimiterConfig{
			Algorithm:       TokenBucketType,
			MaxAttempts:     5,
			Window:          15 * time.Minute,
			LockoutDuration: 15 * time.Minute,
		}

		middleware := NewLoginRateLimiterMiddleware(cfg)

		assert.NotNil(t, middleware)
	})

	t.Run("creates sliding window login middleware", func(t *testing.T) {
		cfg := LoginRateLimiterConfig{
			Algorithm:       SlidingWindowType,
			MaxAttempts:     5,
			Window:          15 * time.Minute,
			LockoutDuration: 15 * time.Minute,
		}

		middleware := NewLoginRateLimiterMiddleware(cfg)

		assert.NotNil(t, middleware)
	})

	t.Run("defaults to token bucket for unknown algorithm", func(t *testing.T) {
		cfg := LoginRateLimiterConfig{
			Algorithm:       "unknown",
			MaxAttempts:     5,
			Window:          15 * time.Minute,
			LockoutDuration: 15 * time.Minute,
		}

		middleware := NewLoginRateLimiterMiddleware(cfg)

		assert.NotNil(t, middleware)
	})
}

func TestNewRateLimiterMiddlewareFromEnv(t *testing.T) {
	t.Run("creates middleware from environment variables", func(t *testing.T) {
		os.Setenv("RATE_LIMIT_ALGORITHM", "token_bucket")
		os.Setenv("RATE_LIMIT_RPM", "150")
		os.Setenv("RATE_LIMIT_BURST", "25")
		defer func() {
			os.Unsetenv("RATE_LIMIT_ALGORITHM")
			os.Unsetenv("RATE_LIMIT_RPM")
			os.Unsetenv("RATE_LIMIT_BURST")
		}()

		middleware := NewRateLimiterMiddlewareFromEnv()

		assert.NotNil(t, middleware)
	})

	t.Run("uses defaults when no env vars set", func(t *testing.T) {
		middleware := NewRateLimiterMiddlewareFromEnv()

		assert.NotNil(t, middleware)
	})
}

func TestNewLoginRateLimiterMiddlewareFromEnv(t *testing.T) {
	t.Run("creates middleware from environment variables", func(t *testing.T) {
		os.Setenv("LOGIN_RATE_LIMIT_ALGORITHM", "token_bucket")
		os.Setenv("LOGIN_MAX_ATTEMPTS", "3")
		os.Setenv("LOGIN_WINDOW_MINUTES", "10")
		os.Setenv("LOGIN_LOCKOUT_MINUTES", "20")
		defer func() {
			os.Unsetenv("LOGIN_RATE_LIMIT_ALGORITHM")
			os.Unsetenv("LOGIN_MAX_ATTEMPTS")
			os.Unsetenv("LOGIN_WINDOW_MINUTES")
			os.Unsetenv("LOGIN_LOCKOUT_MINUTES")
		}()

		middleware := NewLoginRateLimiterMiddlewareFromEnv()

		assert.NotNil(t, middleware)
	})

	t.Run("uses defaults when no env vars set", func(t *testing.T) {
		middleware := NewLoginRateLimiterMiddlewareFromEnv()

		assert.NotNil(t, middleware)
	})
}
