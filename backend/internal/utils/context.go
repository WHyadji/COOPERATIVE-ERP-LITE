package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetKoperasiID safely retrieves cooperative ID from gin context
// Returns error response if cooperative ID is missing or invalid
func GetKoperasiID(c *gin.Context) (uuid.UUID, bool) {
	// Get value from context
	idKoperasi, exists := c.Get("idKoperasi")
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "Cooperative ID not found in context", nil)
		return uuid.Nil, false
	}

	// Type assertion with safety check
	koperasiUUID, ok := idKoperasi.(uuid.UUID)
	if !ok {
		ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Invalid cooperative ID type", nil)
		return uuid.Nil, false
	}

	return koperasiUUID, true
}

// GetPenggunaID safely retrieves user ID from gin context
// Returns error response if user ID is missing or invalid
func GetPenggunaID(c *gin.Context) (uuid.UUID, bool) {
	// Get value from context
	idPengguna, exists := c.Get("idPengguna")
	if !exists {
		ErrorResponse(c, http.StatusUnauthorized, "UNAUTHORIZED", "User ID not found in context", nil)
		return uuid.Nil, false
	}

	// Type assertion with safety check
	penggunaUUID, ok := idPengguna.(uuid.UUID)
	if !ok {
		ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Invalid user ID type", nil)
		return uuid.Nil, false
	}

	return penggunaUUID, true
}

// MustGetKoperasiID retrieves cooperative ID and panics if not found
// Use this only when you're 100% sure middleware has set the value
// (not recommended for production code)
func MustGetKoperasiID(c *gin.Context) uuid.UUID {
	idKoperasi, _ := c.Get("idKoperasi")
	return idKoperasi.(uuid.UUID)
}
