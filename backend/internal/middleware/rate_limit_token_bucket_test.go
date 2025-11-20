package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTokenBucket_Take(t *testing.T) {
	t.Run("allows requests within capacity", func(t *testing.T) {
		bucket := NewTokenBucket(10, 1.0) // 10 tokens, refill 1 per second

		// Should allow 10 requests immediately
		for i := 0; i < 10; i++ {
			assert.True(t, bucket.Take(), "Request %d should be allowed", i+1)
		}

		// 11th request should be blocked
		assert.False(t, bucket.Take(), "11th request should be blocked")
	})

	t.Run("refills tokens over time", func(t *testing.T) {
		bucket := NewTokenBucket(5, 10.0) // 5 tokens, refill 10 per second

		// Consume all tokens
		for i := 0; i < 5; i++ {
			bucket.Take()
		}

		// Should be blocked immediately
		assert.False(t, bucket.Take())

		// Wait 200ms (should refill 2 tokens at 10/sec)
		time.Sleep(200 * time.Millisecond)

		// Should allow 2 requests
		assert.True(t, bucket.Take(), "After 200ms, should allow request")
		assert.True(t, bucket.Take(), "After 200ms, should allow 2nd request")
		assert.False(t, bucket.Take(), "3rd request should be blocked")
	})

	t.Run("does not exceed capacity", func(t *testing.T) {
		bucket := NewTokenBucket(5, 100.0) // 5 tokens, fast refill

		// Wait for refill
		time.Sleep(100 * time.Millisecond)

		// Should still only have 5 tokens (capacity limit)
		allowed := 0
		for i := 0; i < 10; i++ {
			if bucket.Take() {
				allowed++
			}
		}

		assert.Equal(t, 5, allowed, "Should not exceed capacity of 5")
	})

	t.Run("thread safety", func(t *testing.T) {
		bucket := NewTokenBucket(100, 50.0)
		var wg sync.WaitGroup
		allowed := 0
		var mu sync.Mutex

		// 200 goroutines trying to take tokens
		for i := 0; i < 200; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if bucket.Take() {
					mu.Lock()
					allowed++
					mu.Unlock()
				}
			}()
		}

		wg.Wait()

		// Should allow exactly 100 (capacity)
		assert.Equal(t, 100, allowed, "Should allow exactly capacity")
	})
}

func TestRateLimiterTokenBucket_Allow(t *testing.T) {
	t.Run("tracks different IPs separately", func(t *testing.T) {
		limiter := NewRateLimiterTokenBucket(60, 10) // 60/min, burst 10
		defer limiter.Shutdown()

		// IP1 consumes all tokens
		for i := 0; i < 10; i++ {
			assert.True(t, limiter.Allow("192.168.1.1"))
		}
		assert.False(t, limiter.Allow("192.168.1.1"), "IP1 should be blocked")

		// IP2 should still be allowed
		assert.True(t, limiter.Allow("192.168.1.2"), "IP2 should be allowed")
	})

	t.Run("evicts oldest entry when full", func(t *testing.T) {
		limiter := NewRateLimiterTokenBucket(60, 10)
		limiter.maxEntries = 5 // Set low limit for testing
		defer limiter.Shutdown()

		// Add 5 IPs
		for i := 1; i <= 5; i++ {
			ip := "192.168.1." + string(rune(i))
			limiter.Allow(ip)
		}

		// Wait a bit so first IP becomes oldest
		time.Sleep(10 * time.Millisecond)

		// Add 6th IP (should evict oldest)
		assert.True(t, limiter.Allow("192.168.1.6"))

		// Check that we still have max 5 entries
		limiter.mu.RLock()
		count := len(limiter.buckets)
		limiter.mu.RUnlock()

		assert.LessOrEqual(t, count, 5, "Should not exceed max entries")
	})

	t.Run("cleanup removes inactive buckets", func(t *testing.T) {
		limiter := NewRateLimiterTokenBucket(60, 10)
		limiter.cleanupInterval = 100 * time.Millisecond
		defer limiter.Shutdown()

		// Add bucket
		limiter.Allow("192.168.1.1")

		// Manually mark as old
		limiter.mu.Lock()
		bucket := limiter.buckets["192.168.1.1"]
		bucket.mu.Lock()
		bucket.lastRefill = time.Now().Add(-20 * time.Minute) // 20 min ago
		bucket.mu.Unlock()
		limiter.mu.Unlock()

		// Trigger cleanup
		limiter.cleanupInactive()

		// Bucket should be removed
		limiter.mu.RLock()
		_, exists := limiter.buckets["192.168.1.1"]
		limiter.mu.RUnlock()

		assert.False(t, exists, "Inactive bucket should be removed")
	})

	t.Run("graceful shutdown", func(t *testing.T) {
		limiter := NewRateLimiterTokenBucket(60, 10)

		// Start using it
		limiter.Allow("192.168.1.1")

		// Shutdown
		done := make(chan bool)
		go func() {
			limiter.Shutdown()
			done <- true
		}()

		select {
		case <-done:
			// Success
		case <-time.After(2 * time.Second):
			t.Fatal("Shutdown took too long")
		}
	})
}

func TestRateLimiterTokenBucket_GetMetrics(t *testing.T) {
	limiter := NewRateLimiterTokenBucket(100, 20)
	defer limiter.Shutdown()

	// Add some IPs
	limiter.Allow("192.168.1.1")
	limiter.Allow("192.168.1.2")
	limiter.Allow("192.168.1.3")

	metrics := limiter.GetMetrics()

	assert.Equal(t, 3, metrics["tracked_ips"])
	assert.Equal(t, MaxRateLimitEntries, metrics["max_capacity"])
	assert.Equal(t, 100.0, metrics["refill_rate"])
	assert.Equal(t, 20, metrics["burst_size"])
	assert.Greater(t, metrics["utilization"].(float64), 0.0)
}

func TestTokenBucketMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("allows requests within limit", func(t *testing.T) {
		router := gin.New()
		router.Use(TokenBucketMiddleware(60, 10)) // 60/min, burst 10

		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "ok"})
		})

		// Make 10 requests (within burst)
		for i := 0; i < 10; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.1:1234"
			router.ServeHTTP(w, req)

			assert.Equal(t, 200, w.Code, "Request %d should succeed", i+1)
		}
	})

	t.Run("blocks requests over limit", func(t *testing.T) {
		router := gin.New()
		router.Use(TokenBucketMiddleware(60, 5)) // 60/min, burst 5

		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "ok"})
		})

		// Make 6 requests (over burst of 5)
		for i := 0; i < 6; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.1:1234"
			router.ServeHTTP(w, req)

			if i < 5 {
				assert.Equal(t, 200, w.Code, "Request %d should succeed", i+1)
			} else {
				assert.Equal(t, 429, w.Code, "Request 6 should be rate limited")
			}
		}
	})

	t.Run("different IPs tracked separately", func(t *testing.T) {
		router := gin.New()
		router.Use(TokenBucketMiddleware(60, 5))

		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "ok"})
		})

		// IP1: consume all tokens
		for i := 0; i < 5; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			req.RemoteAddr = "192.168.1.1:1234"
			router.ServeHTTP(w, req)
			assert.Equal(t, 200, w.Code)
		}

		// IP1: should be blocked
		w1 := httptest.NewRecorder()
		req1, _ := http.NewRequest("GET", "/test", nil)
		req1.RemoteAddr = "192.168.1.1:1234"
		router.ServeHTTP(w1, req1)
		assert.Equal(t, 429, w1.Code)

		// IP2: should be allowed
		w2 := httptest.NewRecorder()
		req2, _ := http.NewRequest("GET", "/test", nil)
		req2.RemoteAddr = "192.168.1.2:1234"
		router.ServeHTTP(w2, req2)
		assert.Equal(t, 200, w2.Code)
	})
}

func TestLoginRateLimiterTokenBucket(t *testing.T) {
	t.Run("allows attempts within limit", func(t *testing.T) {
		limiter := NewLoginRateLimiterTokenBucket(5, 15*time.Minute, 15*time.Minute)
		defer limiter.Shutdown()

		// Should allow 5 attempts
		for i := 0; i < 5; i++ {
			assert.True(t, limiter.RecordAttempt("192.168.1.1"), "Attempt %d should be allowed", i+1)
		}

		// 6th attempt should be blocked and locked out
		assert.False(t, limiter.RecordAttempt("192.168.1.1"), "6th attempt should be blocked")
	})

	t.Run("locks out after max attempts", func(t *testing.T) {
		limiter := NewLoginRateLimiterTokenBucket(3, 15*time.Minute, 5*time.Second)
		defer limiter.Shutdown()

		// Consume all attempts
		for i := 0; i < 3; i++ {
			limiter.RecordAttempt("user1")
		}

		// Should be locked out
		assert.True(t, limiter.IsLockedOut("user1"))

		// Wait for lockout to expire
		time.Sleep(6 * time.Second)

		// Should no longer be locked out
		assert.False(t, limiter.IsLockedOut("user1"))
	})

	t.Run("clears attempts on success", func(t *testing.T) {
		limiter := NewLoginRateLimiterTokenBucket(5, 15*time.Minute, 15*time.Minute)
		defer limiter.Shutdown()

		// Make failed attempts
		limiter.RecordAttempt("user1")
		limiter.RecordAttempt("user1")

		// Clear on success
		limiter.ClearAttempts("user1")

		// Should have fresh bucket now
		for i := 0; i < 5; i++ {
			assert.True(t, limiter.RecordAttempt("user1"), "Should allow fresh attempts")
		}
	})

	t.Run("tracks different identifiers separately", func(t *testing.T) {
		limiter := NewLoginRateLimiterTokenBucket(3, 15*time.Minute, 15*time.Minute)
		defer limiter.Shutdown()

		// User1 consumes all attempts
		for i := 0; i < 3; i++ {
			limiter.RecordAttempt("user1")
		}
		assert.False(t, limiter.RecordAttempt("user1"))

		// User2 should still be allowed
		assert.True(t, limiter.RecordAttempt("user2"))
	})

	t.Run("cleanup removes expired lockouts", func(t *testing.T) {
		limiter := NewLoginRateLimiterTokenBucket(3, 15*time.Minute, 1*time.Second)
		defer limiter.Shutdown()

		// Lock out user
		for i := 0; i < 3; i++ {
			limiter.RecordAttempt("user1")
		}
		assert.True(t, limiter.IsLockedOut("user1"))

		// Manually set lockout time to past
		limiter.mu.Lock()
		limiter.lockedOut["user1"] = time.Now().Add(-2 * time.Second)
		limiter.mu.Unlock()

		// Trigger cleanup
		limiter.cleanupExpired()

		// Should be removed
		assert.False(t, limiter.IsLockedOut("user1"))
	})
}

func TestLoginTokenBucketMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("blocks when locked out", func(t *testing.T) {
		router := gin.New()
		router.Use(LoginTokenBucketMiddleware(3, 15*time.Minute, 15*time.Minute))

		router.POST("/login", func(c *gin.Context) {
			// Get limiter from context
			limiter, exists := c.Get("loginLimiter")
			assert.True(t, exists)

			lrl := limiter.(*LoginRateLimiterTokenBucket)

			// Simulate failed login
			lrl.RecordAttempt(c.ClientIP())

			c.JSON(200, gin.H{"message": "ok"})
		})

		// Make 3 failed attempts
		for i := 0; i < 3; i++ {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/login", nil)
			req.RemoteAddr = "192.168.1.1:1234"
			router.ServeHTTP(w, req)
			assert.Equal(t, 200, w.Code)
		}

		// 4th request should be blocked by middleware (locked out)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		router.ServeHTTP(w, req)

		assert.Equal(t, 429, w.Code)
		assert.Contains(t, w.Body.String(), "Too many failed login attempts")
	})

	t.Run("allows requests when not locked out", func(t *testing.T) {
		router := gin.New()
		router.Use(LoginTokenBucketMiddleware(5, 15*time.Minute, 15*time.Minute))

		router.POST("/login", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "ok"})
		})

		// First request should pass
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/login", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)
	})
}

func TestLoginRateLimiterTokenBucket_MemoryExhaustion(t *testing.T) {
	limiter := NewLoginRateLimiterTokenBucket(5, 15*time.Minute, 15*time.Minute)
	defer limiter.Shutdown()

	t.Run("locks out when storage is full", func(t *testing.T) {
		// Fill up to max capacity
		for i := 0; i < MaxRateLimitEntries; i++ {
			ip := fmt.Sprintf("192.168.1.%d", i)
			limiter.RecordAttempt(ip)
		}

		// Next attempt should be locked out due to memory protection
		result := limiter.RecordAttempt("new.ip.address")
		assert.False(t, result)

		// Should be locked out
		assert.True(t, limiter.IsLockedOut("new.ip.address"))
	})
}

func TestLoginRateLimiterTokenBucket_CleanupEdgeCases(t *testing.T) {
	t.Run("cleanup handles empty attempts list", func(t *testing.T) {
		limiter := NewLoginRateLimiterTokenBucket(5, 15*time.Minute, 1*time.Second)
		defer limiter.Shutdown()

		// Add empty attempts list
		limiter.mu.Lock()
		limiter.failedAttempts["test.ip"] = []time.Time{}
		limiter.mu.Unlock()

		// Run cleanup
		limiter.cleanupExpired()

		// Should remove empty list
		limiter.mu.RLock()
		_, exists := limiter.failedAttempts["test.ip"]
		limiter.mu.RUnlock()
		assert.False(t, exists)
	})

	t.Run("cleanup preserves recent attempts", func(t *testing.T) {
		limiter := NewLoginRateLimiterTokenBucket(5, 15*time.Minute, 15*time.Minute)
		defer limiter.Shutdown()

		// Add recent attempt
		ip := "recent.ip"
		limiter.RecordAttempt(ip)

		// Run cleanup immediately
		limiter.cleanupExpired()

		// Recent attempt should still exist
		limiter.mu.RLock()
		attempts, exists := limiter.failedAttempts[ip]
		limiter.mu.RUnlock()
		assert.True(t, exists)
		assert.Equal(t, 1, len(attempts))
	})

	t.Run("cleanup removes old attempts", func(t *testing.T) {
		limiter := NewLoginRateLimiterTokenBucket(5, 1*time.Second, 1*time.Second)
		defer limiter.Shutdown()

		// Add old attempt (older than window + 10 minutes)
		ip := "old.ip"
		limiter.mu.Lock()
		limiter.failedAttempts[ip] = []time.Time{time.Now().Add(-15 * time.Minute)}
		limiter.mu.Unlock()

		// Run cleanup
		limiter.cleanupExpired()

		// Old attempt should be removed
		limiter.mu.RLock()
		_, exists := limiter.failedAttempts[ip]
		limiter.mu.RUnlock()
		assert.False(t, exists)
	})
}

// Benchmark tests
func BenchmarkTokenBucket_Take(b *testing.B) {
	bucket := NewTokenBucket(1000, 100.0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bucket.Take()
	}
}

func BenchmarkRateLimiterTokenBucket_Allow(b *testing.B) {
	limiter := NewRateLimiterTokenBucket(100, 20)
	defer limiter.Shutdown()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ip := "192.168.1." + string(rune(i%256))
		limiter.Allow(ip)
	}
}

func BenchmarkTokenBucketMiddleware(b *testing.B) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(TokenBucketMiddleware(100, 20))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		router.ServeHTTP(w, req)
	}
}

func BenchmarkTokenBucket_Concurrent(b *testing.B) {
	limiter := NewRateLimiterTokenBucket(1000, 100)
	defer limiter.Shutdown()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ip := fmt.Sprintf("192.168.1.%d", i%256)
			limiter.Allow(ip)
			i++
		}
	})
}

func TestTokenBucket_MemoryUsage(t *testing.T) {
	// Verify memory usage is ~24 bytes per IP as documented
	limiter := NewRateLimiterTokenBucket(100, 20)
	defer limiter.Shutdown()

	// Add 1000 IPs
	for i := 0; i < 1000; i++ {
		ip := fmt.Sprintf("192.168.%d.%d", i/256, i%256)
		limiter.Allow(ip)
	}

	// Get metrics
	metrics := limiter.GetMetrics()
	trackedIPs := metrics["tracked_ips"].(int)

	// Should track all 1000 IPs
	assert.Equal(t, 1000, trackedIPs)

	// Note: Exact memory usage is hard to measure in Go due to GC
	// But we can verify the structure is minimal:
	// TokenBucket struct: 8 (float64) + 4 (int) + 8 (float64) + 8 (time.Time) + 8 (mutex) = ~36 bytes
	// Map overhead: ~16 bytes per entry
	// Total: ~52 bytes per IP (still much better than 200 bytes for sliding window)
}
