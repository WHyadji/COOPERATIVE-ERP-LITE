package handlers

import (
	"cooperative-erp-lite/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ParameterPaginasi holds validated pagination parameters
type ParameterPaginasi struct {
	Halaman       int
	UkuranHalaman int
}

// AmbilIDKoperasiDariContext extracts and validates the cooperative ID from the Gin context.
// This function is used for multi-tenant validation to ensure all operations are scoped
// to the correct cooperative.
//
// Parameters:
//   - c: The Gin context containing the idKoperasi set by the auth middleware
//
// Returns:
//   - uuid.UUID: The cooperative ID from the context
//   - bool: true if successful, false if the ID is missing or invalid
//
// If the ID is missing or invalid, this function sends an appropriate error response
// and returns false. The caller should return immediately when ok is false.
func AmbilIDKoperasiDariContext(c *gin.Context) (uuid.UUID, bool) {
	idKoperasi, exists := c.Get("idKoperasi")
	if !exists {
		utils.InternalServerErrorResponse(c, "ID koperasi tidak ditemukan dalam context", nil)
		return uuid.Nil, false
	}

	koperasiUUID, ok := idKoperasi.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "ID koperasi tidak valid", nil)
		return uuid.Nil, false
	}

	return koperasiUUID, true
}

// AmbilIDPenggunaDariContext extracts and validates the user ID from the Gin context.
// This function is used to identify the user performing an action for audit logging
// and authorization purposes.
//
// Parameters:
//   - c: The Gin context containing the idPengguna set by the auth middleware
//
// Returns:
//   - uuid.UUID: The user ID from the context
//   - bool: true if successful, false if the ID is missing or invalid
//
// If the ID is missing or invalid, this function sends an appropriate error response
// and returns false. The caller should return immediately when ok is false.
func AmbilIDPenggunaDariContext(c *gin.Context) (uuid.UUID, bool) {
	idPengguna, exists := c.Get("idPengguna")
	if !exists {
		utils.UnauthorizedResponse(c, "Token tidak valid")
		return uuid.Nil, false
	}

	penggunaUUID, ok := idPengguna.(uuid.UUID)
	if !ok {
		utils.UnauthorizedResponse(c, "ID pengguna tidak valid")
		return uuid.Nil, false
	}

	return penggunaUUID, true
}

// AmbilParameterPaginasi extracts and validates pagination parameters from the request query.
// This function applies security validation to prevent DoS attacks from excessive page sizes.
//
// Query Parameters:
//   - page: The page number (default: 1, minimum: 1)
//   - pageSize: Items per page (default: 20, minimum: 1, maximum: 100)
//
// Returns:
//   - ParameterPaginasi: Struct containing validated Halaman and UkuranHalaman
//
// Security features:
//   - Validates and clamps page numbers to prevent negative values
//   - Enforces maximum page size of 100 to prevent memory exhaustion DoS
//   - Uses secure defaults when parameters are missing or invalid
func AmbilParameterPaginasi(c *gin.Context) ParameterPaginasi {
	// Parse query parameters with defaults
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))

	// Apply validation using the security-hardened utility function
	validPage, validPageSize := utils.ValidatePagination(page, pageSize)

	return ParameterPaginasi{
		Halaman:       validPage,
		UkuranHalaman: validPageSize,
	}
}

// ParseUUIDDariParameter extracts and parses a UUID from a URL parameter.
// This function validates that the parameter exists and is a valid UUID.
//
// Parameters:
//   - c: The Gin context containing the URL parameters
//   - paramName: The name of the parameter to extract (e.g., "id")
//
// Returns:
//   - uuid.UUID: The parsed UUID value
//   - bool: true if successful, false if the parameter is missing or invalid
//
// If the parameter is missing or invalid, this function sends an appropriate error response
// and returns false. The caller should return immediately when ok is false.
func ParseUUIDDariParameter(c *gin.Context, paramName string) (uuid.UUID, bool) {
	idStr := c.Param(paramName)
	if idStr == "" {
		utils.BadRequestResponse(c, "Parameter "+paramName+" diperlukan")
		return uuid.Nil, false
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.BadRequestResponse(c, "Parameter "+paramName+" bukan UUID yang valid")
		return uuid.Nil, false
	}

	return id, true
}
