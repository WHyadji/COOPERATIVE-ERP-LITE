package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter tracks request rates per IP address
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int           // max requests
	window   time.Duration // time window
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// cleanup removes old entries periodically
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, times := range rl.requests {
			// Remove timestamps older than window
			filtered := make([]time.Time, 0)
			for _, t := range times {
				if now.Sub(t) < rl.window {
					filtered = append(filtered, t)
				}
			}
			if len(filtered) == 0 {
				delete(rl.requests, ip)
			} else {
				rl.requests[ip] = filtered
			}
		}
		rl.mu.Unlock()
	}
}

// Allow checks if a request from the given IP should be allowed
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Get existing timestamps for this IP
	times, exists := rl.requests[ip]
	if !exists {
		rl.requests[ip] = []time.Time{now}
		return true
	}

	// Filter out timestamps outside the window
	filtered := make([]time.Time, 0)
	for _, t := range times {
		if now.Sub(t) < rl.window {
			filtered = append(filtered, t)
		}
	}

	// Check if limit exceeded
	if len(filtered) >= rl.limit {
		return false
	}

	// Add current request
	filtered = append(filtered, now)
	rl.requests[ip] = filtered

	return true
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(limit int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(limit, window)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. Please try again later.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// LoginRateLimiter is a specialized rate limiter for login endpoints
type LoginRateLimiter struct {
	attempts map[string][]time.Time
	mu       sync.RWMutex
	maxAttempts int
	window      time.Duration
	lockoutDuration time.Duration
	lockedOut   map[string]time.Time
}

// NewLoginRateLimiter creates a rate limiter specifically for login attempts
func NewLoginRateLimiter(maxAttempts int, window time.Duration, lockoutDuration time.Duration) *LoginRateLimiter {
	lrl := &LoginRateLimiter{
		attempts: make(map[string][]time.Time),
		lockedOut: make(map[string]time.Time),
		maxAttempts: maxAttempts,
		window: window,
		lockoutDuration: lockoutDuration,
	}

	go lrl.cleanup()

	return lrl
}

// cleanup removes old entries periodically
func (lrl *LoginRateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		lrl.mu.Lock()
		now := time.Now()

		// Clean up attempts
		for key, times := range lrl.attempts {
			filtered := make([]time.Time, 0)
			for _, t := range times {
				if now.Sub(t) < lrl.window {
					filtered = append(filtered, t)
				}
			}
			if len(filtered) == 0 {
				delete(lrl.attempts, key)
			} else {
				lrl.attempts[key] = filtered
			}
		}

		// Clean up lockouts
		for key, lockoutTime := range lrl.lockedOut {
			if now.Sub(lockoutTime) > lrl.lockoutDuration {
				delete(lrl.lockedOut, key)
			}
		}

		lrl.mu.Unlock()
	}
}

// IsLockedOut checks if an IP/username is currently locked out
func (lrl *LoginRateLimiter) IsLockedOut(identifier string) bool {
	lrl.mu.RLock()
	defer lrl.mu.RUnlock()

	lockoutTime, exists := lrl.lockedOut[identifier]
	if !exists {
		return false
	}

	// Check if lockout period has expired
	if time.Since(lockoutTime) > lrl.lockoutDuration {
		return false
	}

	return true
}

// RecordAttempt records a failed login attempt
func (lrl *LoginRateLimiter) RecordAttempt(identifier string) bool {
	lrl.mu.Lock()
	defer lrl.mu.Unlock()

	// Check if already locked out
	if lockoutTime, exists := lrl.lockedOut[identifier]; exists {
		if time.Since(lockoutTime) < lrl.lockoutDuration {
			return false // Still locked out
		}
		// Lockout expired, remove it
		delete(lrl.lockedOut, identifier)
	}

	now := time.Now()
	times, exists := lrl.attempts[identifier]
	if !exists {
		lrl.attempts[identifier] = []time.Time{now}
		return true
	}

	// Filter recent attempts within window
	filtered := make([]time.Time, 0)
	for _, t := range times {
		if now.Sub(t) < lrl.window {
			filtered = append(filtered, t)
		}
	}

	// Add current attempt
	filtered = append(filtered, now)
	lrl.attempts[identifier] = filtered

	// Check if exceeded max attempts
	if len(filtered) >= lrl.maxAttempts {
		lrl.lockedOut[identifier] = now
		return false
	}

	return true
}

// ClearAttempts clears login attempts for an identifier (on successful login)
func (lrl *LoginRateLimiter) ClearAttempts(identifier string) {
	lrl.mu.Lock()
	defer lrl.mu.Unlock()

	delete(lrl.attempts, identifier)
	delete(lrl.lockedOut, identifier)
}

// LoginRateLimitMiddleware creates rate limiting middleware for login endpoints
func LoginRateLimitMiddleware(maxAttempts int, window time.Duration, lockoutDuration time.Duration) gin.HandlerFunc {
	limiter := NewLoginRateLimiter(maxAttempts, window, lockoutDuration)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		// Check if IP is locked out
		if limiter.IsLockedOut(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many failed login attempts. Account temporarily locked.",
			})
			c.Abort()
			return
		}

		// Store limiter in context for use in handler
		c.Set("loginLimiter", limiter)

		c.Next()
	}
}
