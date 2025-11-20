package middleware

import (
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// TokenBucket implements token bucket rate limiting algorithm
// Memory efficient: ~24 bytes per IP vs ~200 bytes for sliding window log
// Time complexity: O(1) per request vs O(n) for sliding window
type TokenBucket struct {
	tokens     float64   // Current available tokens
	capacity   int       // Max tokens (burst size)
	refillRate float64   // Tokens added per second
	lastRefill time.Time // Last time tokens were refilled
	mu         sync.Mutex
}

// NewTokenBucket creates a new token bucket for an IP
func NewTokenBucket(capacity int, refillRate float64) *TokenBucket {
	return &TokenBucket{
		tokens:     float64(capacity), // Start with full bucket
		capacity:   capacity,
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

// Take attempts to take one token from the bucket
// Returns true if token was available, false if rate limited
func (tb *TokenBucket) Take() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()

	// Refill tokens based on elapsed time
	tb.tokens = math.Min(
		float64(tb.capacity),
		tb.tokens+elapsed*tb.refillRate,
	)
	tb.lastRefill = now

	// Try to take 1 token
	if tb.tokens >= 1.0 {
		tb.tokens -= 1.0
		return true
	}

	return false
}

// RateLimiterTokenBucket manages token buckets for multiple IPs
type RateLimiterTokenBucket struct {
	buckets      map[string]*TokenBucket
	mu           sync.RWMutex
	capacity     int
	refillRate   float64
	stopChan     chan struct{}
	wg           sync.WaitGroup
	maxEntries   int
	cleanupInterval time.Duration
}

// NewRateLimiterTokenBucket creates a new token bucket rate limiter
// requestsPerMinute: Maximum requests allowed per minute
// burstSize: Maximum burst allowed (capacity)
func NewRateLimiterTokenBucket(requestsPerMinute int, burstSize int) *RateLimiterTokenBucket {
	rl := &RateLimiterTokenBucket{
		buckets:         make(map[string]*TokenBucket),
		capacity:        burstSize,
		refillRate:      float64(requestsPerMinute) / 60.0, // Convert to per second
		stopChan:        make(chan struct{}),
		maxEntries:      MaxRateLimitEntries, // 10,000 (from rate_limit.go)
		cleanupInterval: 5 * time.Minute,
	}

	// Start cleanup goroutine with graceful shutdown
	rl.wg.Add(1)
	go rl.cleanup()

	return rl
}

// Allow checks if a request from the given IP should be allowed
func (rl *RateLimiterTokenBucket) Allow(ip string) bool {
	rl.mu.Lock()
	bucket, exists := rl.buckets[ip]
	if !exists {
		// Prevent memory exhaustion
		if len(rl.buckets) >= rl.maxEntries {
			// Evict oldest entry instead of rejecting request
			rl.evictOldestUnsafe()
		}

		bucket = NewTokenBucket(rl.capacity, rl.refillRate)
		rl.buckets[ip] = bucket
	}
	rl.mu.Unlock()

	// Check if allowed (without holding global lock)
	return bucket.Take()
}

// evictOldestUnsafe removes the oldest (least recently used) bucket
// Must be called with rl.mu locked
func (rl *RateLimiterTokenBucket) evictOldestUnsafe() {
	var oldestIP string
	var oldestTime time.Time = time.Now()

	for ip, bucket := range rl.buckets {
		bucket.mu.Lock()
		if bucket.lastRefill.Before(oldestTime) {
			oldestTime = bucket.lastRefill
			oldestIP = ip
		}
		bucket.mu.Unlock()
	}

	if oldestIP != "" {
		delete(rl.buckets, oldestIP)
	}
}

// cleanup removes inactive buckets periodically
func (rl *RateLimiterTokenBucket) cleanup() {
	defer rl.wg.Done()
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.cleanupInactive()
		case <-rl.stopChan:
			return // Graceful shutdown
		}
	}
}

// cleanupInactive removes buckets that haven't been used in 10 minutes
func (rl *RateLimiterTokenBucket) cleanupInactive() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	inactiveThreshold := 10 * time.Minute

	for ip, bucket := range rl.buckets {
		bucket.mu.Lock()
		inactive := now.Sub(bucket.lastRefill) > inactiveThreshold
		bucket.mu.Unlock()

		if inactive {
			delete(rl.buckets, ip)
		}
	}
}

// Shutdown stops the cleanup goroutine gracefully
func (rl *RateLimiterTokenBucket) Shutdown() {
	close(rl.stopChan)
	rl.wg.Wait()
}

// GetMetrics returns current rate limiter metrics
func (rl *RateLimiterTokenBucket) GetMetrics() map[string]interface{} {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	return map[string]interface{}{
		"tracked_ips":   len(rl.buckets),
		"max_capacity":  rl.maxEntries,
		"utilization":   float64(len(rl.buckets)) / float64(rl.maxEntries) * 100,
		"refill_rate":   rl.refillRate * 60, // Convert to requests per minute for readability
		"burst_size":    rl.capacity,
	}
}

// TokenBucketMiddleware creates a rate limiting middleware using token bucket algorithm
// requestsPerMinute: Maximum requests allowed per minute (e.g., 100)
// burstSize: Maximum burst allowed (e.g., 20 for 20 simultaneous requests)
func TokenBucketMiddleware(requestsPerMinute int, burstSize int) gin.HandlerFunc {
	limiter := NewRateLimiterTokenBucket(requestsPerMinute, burstSize)

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

// LoginRateLimiterTokenBucket is a specialized rate limiter for login endpoints
// It tracks failed login attempts separately from rate limiting
type LoginRateLimiterTokenBucket struct {
	failedAttempts  map[string][]time.Time // Track failed attempt timestamps
	lockedOut       map[string]time.Time   // Track lockout start times
	mu              sync.RWMutex
	maxAttempts     int           // Max failed attempts before lockout
	window          time.Duration // Time window for counting attempts
	lockoutDuration time.Duration // How long to lock out
	stopChan        chan struct{}
	wg              sync.WaitGroup
}

// NewLoginRateLimiterTokenBucket creates a login rate limiter
// It tracks failed attempts and enforces lockouts after max attempts
func NewLoginRateLimiterTokenBucket(maxAttempts int, window time.Duration, lockoutDuration time.Duration) *LoginRateLimiterTokenBucket {
	lrl := &LoginRateLimiterTokenBucket{
		failedAttempts:  make(map[string][]time.Time),
		lockedOut:       make(map[string]time.Time),
		maxAttempts:     maxAttempts,
		window:          window,
		lockoutDuration: lockoutDuration,
		stopChan:        make(chan struct{}),
	}

	lrl.wg.Add(1)
	go lrl.cleanup()

	return lrl
}

// IsLockedOut checks if an identifier is currently locked out
func (lrl *LoginRateLimiterTokenBucket) IsLockedOut(identifier string) bool {
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
// Returns true if attempt is allowed, false if locked out
func (lrl *LoginRateLimiterTokenBucket) RecordAttempt(identifier string) bool {
	lrl.mu.Lock()
	defer lrl.mu.Unlock()

	// Check if already locked out
	if lockoutTime, exists := lrl.lockedOut[identifier]; exists {
		if time.Since(lockoutTime) < lrl.lockoutDuration {
			return false // Still locked out
		}
		// Lockout expired, remove it
		delete(lrl.lockedOut, identifier)
		delete(lrl.failedAttempts, identifier)
	}

	now := time.Now()

	// Get existing failed attempts
	attempts, exists := lrl.failedAttempts[identifier]
	if !exists {
		// Prevent memory exhaustion
		if len(lrl.failedAttempts) >= MaxRateLimitEntries {
			// Lock out when storage is full to prevent DoS
			lrl.lockedOut[identifier] = now
			return false
		}
		attempts = []time.Time{}
	}

	// Filter attempts within the time window
	validAttempts := make([]time.Time, 0)
	for _, attemptTime := range attempts {
		if now.Sub(attemptTime) < lrl.window {
			validAttempts = append(validAttempts, attemptTime)
		}
	}

	// Add current failed attempt
	validAttempts = append(validAttempts, now)
	lrl.failedAttempts[identifier] = validAttempts

	// Check if reached or exceeded max attempts
	// Lock out when attempts >= maxAttempts
	// Example: maxAttempts=5 means:
	//   - Attempts 1-5: Allowed (return true), lock out after 5th
	//   - Attempt 6+: Blocked by lockout check at top (return false)
	if len(validAttempts) >= lrl.maxAttempts {
		lrl.lockedOut[identifier] = now
	}

	return true // This attempt is allowed (lockout applies to NEXT attempt)
}

// ClearAttempts clears login attempts for an identifier (on successful login)
func (lrl *LoginRateLimiterTokenBucket) ClearAttempts(identifier string) {
	lrl.mu.Lock()
	defer lrl.mu.Unlock()

	delete(lrl.failedAttempts, identifier)
	delete(lrl.lockedOut, identifier)
}

// cleanup removes old entries periodically
func (lrl *LoginRateLimiterTokenBucket) cleanup() {
	defer lrl.wg.Done()
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			lrl.cleanupExpired()
		case <-lrl.stopChan:
			return
		}
	}
}

// cleanupExpired removes expired lockouts and old failed attempts
func (lrl *LoginRateLimiterTokenBucket) cleanupExpired() {
	lrl.mu.Lock()
	defer lrl.mu.Unlock()

	now := time.Now()

	// Clean up expired lockouts
	for identifier, lockoutTime := range lrl.lockedOut {
		if now.Sub(lockoutTime) > lrl.lockoutDuration {
			delete(lrl.lockedOut, identifier)
		}
	}

	// Clean up old failed attempts (older than window + 10 minutes)
	cleanupThreshold := lrl.window + 10*time.Minute
	for identifier, attempts := range lrl.failedAttempts {
		if len(attempts) == 0 {
			delete(lrl.failedAttempts, identifier)
			continue
		}

		// Check if all attempts are old
		allOld := true
		for _, attemptTime := range attempts {
			if now.Sub(attemptTime) < cleanupThreshold {
				allOld = false
				break
			}
		}

		if allOld {
			delete(lrl.failedAttempts, identifier)
		}
	}
}

// Shutdown stops the cleanup goroutine gracefully
func (lrl *LoginRateLimiterTokenBucket) Shutdown() {
	close(lrl.stopChan)
	lrl.wg.Wait()
}

// LoginTokenBucketMiddleware creates rate limiting middleware for login endpoints
func LoginTokenBucketMiddleware(maxAttempts int, window time.Duration, lockoutDuration time.Duration) gin.HandlerFunc {
	limiter := NewLoginRateLimiterTokenBucket(maxAttempts, window, lockoutDuration)

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
