package utils

import (
	"context"
	"time"
)

// Pagination constants for security and performance
const (
	// DefaultPageSize is the default number of items per page when not specified
	DefaultPageSize = 20

	// MaxPageSize is the maximum allowed page size to prevent DoS attacks
	MaxPageSize = 100

	// MinPageSize is the minimum allowed page size
	MinPageSize = 1

	// MinPage is the minimum page number (1-indexed pagination)
	MinPage = 1

	// QueryTimeout is the maximum duration for database queries
	QueryTimeout = 30 * time.Second
)

// PaginationMeta contains pagination metadata for API responses
type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
	HasNext    bool  `json:"hasNext"`
	HasPrev    bool  `json:"hasPrev"`
}

// ValidatePagination validates and normalizes pagination parameters
//
// Parameters:
//   - page: The requested page number (1-indexed)
//   - pageSize: The requested number of items per page
//
// Returns:
//   - validPage: Normalized page number (minimum 1)
//   - validPageSize: Normalized page size (clamped to MinPageSize-MaxPageSize range)
//
// Security features:
//   - Prevents negative page numbers (DoS risk)
//   - Clamps page size to MaxPageSize (memory exhaustion risk)
//   - Ensures minimum valid values
//   - Uses DefaultPageSize for invalid inputs
func ValidatePagination(page, pageSize int) (validPage int, validPageSize int) {
	// Validate and normalize page number
	// Prevent negative pages which could cause invalid SQL offsets
	if page < MinPage {
		validPage = MinPage
	} else {
		validPage = page
	}

	// Validate and normalize page size
	// Prevent DoS attacks from requesting unlimited records
	if pageSize < MinPageSize || pageSize > MaxPageSize {
		// Use default for out-of-range values
		if pageSize < MinPageSize {
			validPageSize = DefaultPageSize
		} else {
			// Clamp to maximum to prevent memory exhaustion
			validPageSize = MaxPageSize
		}
	} else {
		validPageSize = pageSize
	}

	return validPage, validPageSize
}

// CalculatePaginationMeta calculates pagination metadata for API responses
//
// Parameters:
//   - page: Current page number (validated)
//   - pageSize: Number of items per page (validated)
//   - total: Total number of items in the dataset
//
// Returns:
//   - PaginationMeta with all calculated fields
func CalculatePaginationMeta(page, pageSize int, total int64) *PaginationMeta {
	// Calculate total pages (ceiling division)
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	// Prevent division by zero
	if totalPages == 0 {
		totalPages = 1
	}

	// Calculate navigation flags
	hasNext := page < totalPages
	hasPrev := page > MinPage

	return &PaginationMeta{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    hasNext,
		HasPrev:    hasPrev,
	}
}

// CreateQueryContext creates a context with timeout for database queries
//
// This prevents long-running queries from exhausting database resources
// and provides a mechanism to cancel expensive operations.
//
// Returns:
//   - context.Context with timeout set to QueryTimeout
//   - cancel function that must be called to release resources
//
// Example usage:
//
//	ctx, cancel := CreateQueryContext()
//	defer cancel()
//	db.WithContext(ctx).Find(&results)
func CreateQueryContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), QueryTimeout)
}

// CalculateOffset calculates the database offset for pagination
//
// Parameters:
//   - page: Current page number (1-indexed)
//   - pageSize: Number of items per page
//
// Returns:
//   - offset: Number of records to skip in database query
//
// Formula: offset = (page - 1) * pageSize
func CalculateOffset(page, pageSize int) int {
	if page < MinPage {
		page = MinPage
	}
	return (page - 1) * pageSize
}
