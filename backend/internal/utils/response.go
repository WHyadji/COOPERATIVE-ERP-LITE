package utils

import (
	"github.com/gin-gonic/gin"
)

// APIResponse adalah struktur standar untuk response API
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

// ErrorDetail adalah detail error untuk response
type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// PaginationMeta adalah metadata untuk pagination
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalPages int   `json:"totalPages"`
	TotalItems int64 `json:"totalItems"`
}

// PaginatedResponse adalah response dengan pagination
type PaginatedResponse struct {
	Success    bool           `json:"success"`
	Message    string         `json:"message"`
	Data       interface{}    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// SuccessResponse mengirim response sukses
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse mengirim response error
func ErrorResponse(c *gin.Context, statusCode int, code string, message string, details interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: "Terjadi kesalahan",
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// PaginatedSuccessResponse mengirim response sukses dengan pagination
func PaginatedSuccessResponse(c *gin.Context, statusCode int, message string, data interface{}, pagination PaginationMeta) {
	c.JSON(statusCode, PaginatedResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}

// ValidationErrorResponse mengirim response untuk validation error
// DEPRECATED: Use SafeValidationErrorResponse for error objects to prevent information disclosure
func ValidationErrorResponse(c *gin.Context, errors interface{}) {
	ErrorResponse(c, 400, "VALIDATION_ERROR", "Validasi gagal", errors)
}

// SafeValidationErrorResponse mengirim response untuk validation error dengan sanitasi otomatis
// Gunakan fungsi ini ketika menerima error object untuk mencegah information disclosure
func SafeValidationErrorResponse(c *gin.Context, err error) {
	sanitizedMsg := SanitizeError(err)
	ErrorResponse(c, 400, "VALIDATION_ERROR", sanitizedMsg, nil)
}

// BadRequestResponse mengirim response untuk bad request
func BadRequestResponse(c *gin.Context, message string) {
	ErrorResponse(c, 400, "BAD_REQUEST", message, nil)
}

// UnauthorizedResponse mengirim response untuk unauthorized access
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, 401, "UNAUTHORIZED", message, nil)
}

// ForbiddenResponse mengirim response untuk forbidden access
func ForbiddenResponse(c *gin.Context, message string) {
	ErrorResponse(c, 403, "FORBIDDEN", message, nil)
}

// NotFoundResponse mengirim response untuk resource not found
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, 404, "NOT_FOUND", message, nil)
}

// ConflictResponse mengirim response untuk conflict
func ConflictResponse(c *gin.Context, message string) {
	ErrorResponse(c, 409, "CONFLICT", message, nil)
}

// InternalServerErrorResponse mengirim response untuk internal server error
// DEPRECATED: Use SafeInternalServerErrorResponse for error objects to prevent information disclosure
func InternalServerErrorResponse(c *gin.Context, message string, details interface{}) {
	ErrorResponse(c, 500, "INTERNAL_SERVER_ERROR", message, details)
}

// SafeInternalServerErrorResponse mengirim response untuk internal server error dengan sanitasi otomatis
// Gunakan fungsi ini ketika menerima error object untuk mencegah information disclosure
// Error details akan di-log secara internal tetapi tidak dikirim ke client
func SafeInternalServerErrorResponse(c *gin.Context, err error) {
	// TODO: Log the actual error internally for debugging when logger is implemented
	// For now, error details are sanitized and only safe messages are sent to client

	// Send sanitized message to client
	sanitizedMsg := SanitizeError(err)
	ErrorResponse(c, 500, "INTERNAL_SERVER_ERROR", sanitizedMsg, nil)
}

// CalculatePaginationMeta menghitung metadata pagination
func CalculatePaginationMeta(page, pageSize int, totalItems int64) PaginationMeta {
	totalPages := int(totalItems) / pageSize
	if int(totalItems)%pageSize > 0 {
		totalPages++
	}

	return PaginationMeta{
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
		TotalItems: totalItems,
	}
}
