// ============================================================================
// Row-Level Security (RLS) Middleware
// Sets cooperative context for database queries
// ============================================================================

package middleware

import (
	"cooperative-erp-lite/internal/config"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RLSMiddleware sets the current cooperative ID in the database session
// This enables PostgreSQL Row-Level Security policies to filter data automatically
//
// IMPORTANT: This middleware MUST run after AuthMiddleware
// because it depends on the cooperative_id from the JWT token
func RLSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get cooperative ID from context (set by AuthMiddleware or AuthAnggotaMiddleware)
		cooperativeID, exists := c.Get("idKoperasi")
		if !exists {
			// Skip RLS for public endpoints (like health check, login)
			c.Next()
			return
		}

		// Convert to UUID (handles both string and uuid.UUID types)
		var coopUUID uuid.UUID
		var err error

		switch v := cooperativeID.(type) {
		case string:
			coopUUID, err = uuid.Parse(v)
			if err != nil {
				log.Printf("RLS Middleware: Invalid cooperative ID format: %v", err)
				c.Next()
				return
			}
		case uuid.UUID:
			coopUUID = v
		default:
			log.Printf("RLS Middleware: Unexpected cooperative ID type: %T", cooperativeID)
			c.Next()
			return
		}

		// Set RLS context in database
		if err := SetCooperativeContext(coopUUID); err != nil {
			log.Printf("RLS Middleware: Failed to set cooperative context: %v", err)
			// Continue anyway - application-level filtering will still work
		} else {
			log.Printf("RLS: Set cooperative context to %s", coopUUID)
		}

		c.Next()
	}
}

// SetCooperativeContext sets the current cooperative ID in the database session
// This is called by the middleware to enable Row-Level Security filtering
func SetCooperativeContext(cooperativeID uuid.UUID) error {
	db := config.GetDB()
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	// Execute the PostgreSQL function to set session variable
	query := "SELECT set_current_koperasi_id($1)"
	result := db.Exec(query, cooperativeID)

	if result.Error != nil {
		return fmt.Errorf("failed to set cooperative context: %w", result.Error)
	}

	return nil
}

// ClearCooperativeContext clears the current cooperative ID from the database session
// This should be called when logging out or resetting the session
func ClearCooperativeContext() error {
	db := config.GetDB()
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	// Execute the PostgreSQL function to clear session variable
	query := "SELECT clear_current_koperasi_id()"
	result := db.Exec(query)

	if result.Error != nil {
		return fmt.Errorf("failed to clear cooperative context: %w", result.Error)
	}

	return nil
}

// GetCurrentCooperativeContext retrieves the current cooperative ID from the database session
// Useful for debugging and verification
func GetCurrentCooperativeContext() (uuid.UUID, error) {
	db := config.GetDB()
	if db == nil {
		return uuid.Nil, fmt.Errorf("database connection is nil")
	}

	var cooperativeID uuid.UUID
	query := "SELECT get_current_koperasi_id()"
	result := db.Raw(query).Scan(&cooperativeID)

	if result.Error != nil {
		return uuid.Nil, fmt.Errorf("failed to get cooperative context: %w", result.Error)
	}

	return cooperativeID, nil
}
