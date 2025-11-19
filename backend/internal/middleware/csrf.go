package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	CSRFTokenHeader = "X-CSRF-Token"
	CSRFTokenCookie = "csrf_token"
	CSRFTokenLength = 32
)

// CSRFStore stores CSRF tokens with expiration
type CSRFStore struct {
	tokens map[string]time.Time
	mu     sync.RWMutex
}

var csrfStore = &CSRFStore{
	tokens: make(map[string]time.Time),
}

// GenerateCSRFToken generates a new CSRF token
func GenerateCSRFToken() (string, error) {
	b := make([]byte, CSRFTokenLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(b)

	// Store token with 24 hour expiration
	csrfStore.mu.Lock()
	csrfStore.tokens[token] = time.Now().Add(24 * time.Hour)
	csrfStore.mu.Unlock()

	// Clean up expired tokens
	go csrfStore.cleanExpired()

	return token, nil
}

// cleanExpired removes expired tokens from the store
func (s *CSRFStore) cleanExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for token, expiry := range s.tokens {
		if now.After(expiry) {
			delete(s.tokens, token)
		}
	}
}

// ValidateCSRFToken validates a CSRF token
func ValidateCSRFToken(token string) bool {
	csrfStore.mu.RLock()
	defer csrfStore.mu.RUnlock()

	expiry, exists := csrfStore.tokens[token]
	if !exists {
		return false
	}

	return time.Now().Before(expiry)
}

// CSRFProtection middleware validates CSRF tokens for state-changing requests
func CSRFProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip CSRF check for safe methods (GET, HEAD, OPTIONS)
		if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Get token from header
		token := c.GetHeader(CSRFTokenHeader)
		if token == "" {
			// Also check form field as fallback
			token = c.PostForm("csrf_token")
		}

		if token == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "CSRF token missing",
			})
			c.Abort()
			return
		}

		// Validate token
		if !ValidateCSRFToken(token) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Invalid or expired CSRF token",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GenerateCSRFTokenEndpoint returns a new CSRF token to clients
func GenerateCSRFTokenEndpoint(c *gin.Context) {
	token, err := GenerateCSRFToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate CSRF token",
		})
		return
	}

	// Set as cookie for browser clients
	c.SetCookie(
		CSRFTokenCookie,
		token,
		86400, // 24 hours
		"/",
		"",
		false, // secure (set to true in production with HTTPS)
		true,  // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
