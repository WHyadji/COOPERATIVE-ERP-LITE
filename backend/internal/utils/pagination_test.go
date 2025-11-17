package utils

import (
	"context"
	"testing"
	"time"
)

func TestValidatePagination(t *testing.T) {
	tests := []struct {
		name             string
		inputPage        int
		inputPageSize    int
		expectedPage     int
		expectedPageSize int
		description      string
	}{
		{
			name:             "Valid pagination parameters",
			inputPage:        1,
			inputPageSize:    20,
			expectedPage:     1,
			expectedPageSize: 20,
			description:      "Should accept valid page and pageSize",
		},
		{
			name:             "Negative page number (DoS risk)",
			inputPage:        -1,
			inputPageSize:    20,
			expectedPage:     1,
			expectedPageSize: 20,
			description:      "Should normalize negative page to MinPage",
		},
		{
			name:             "Zero page number",
			inputPage:        0,
			inputPageSize:    20,
			expectedPage:     1,
			expectedPageSize: 20,
			description:      "Should normalize zero page to MinPage",
		},
		{
			name:             "Page size exceeds maximum (DoS risk)",
			inputPage:        1,
			inputPageSize:    999999999,
			expectedPage:     1,
			expectedPageSize: MaxPageSize,
			description:      "Should clamp pageSize to MaxPageSize to prevent memory exhaustion",
		},
		{
			name:             "Negative page size",
			inputPage:        1,
			inputPageSize:    -10,
			expectedPage:     1,
			expectedPageSize: DefaultPageSize,
			description:      "Should use DefaultPageSize for negative pageSize",
		},
		{
			name:             "Zero page size",
			inputPage:        1,
			inputPageSize:    0,
			expectedPage:     1,
			expectedPageSize: DefaultPageSize,
			description:      "Should use DefaultPageSize for zero pageSize",
		},
		{
			name:             "Both parameters invalid",
			inputPage:        -5,
			inputPageSize:    -10,
			expectedPage:     1,
			expectedPageSize: DefaultPageSize,
			description:      "Should normalize both parameters when invalid",
		},
		{
			name:             "Maximum valid page size",
			inputPage:        1,
			inputPageSize:    MaxPageSize,
			expectedPage:     1,
			expectedPageSize: MaxPageSize,
			description:      "Should accept MaxPageSize as valid",
		},
		{
			name:             "Minimum valid page size",
			inputPage:        1,
			inputPageSize:    MinPageSize,
			expectedPage:     1,
			expectedPageSize: MinPageSize,
			description:      "Should accept MinPageSize as valid",
		},
		{
			name:             "Large page number",
			inputPage:        1000000,
			inputPageSize:    20,
			expectedPage:     1000000,
			expectedPageSize: 20,
			description:      "Should accept large but valid page numbers",
		},
		{
			name:             "Page size just over maximum",
			inputPage:        1,
			inputPageSize:    MaxPageSize + 1,
			expectedPage:     1,
			expectedPageSize: MaxPageSize,
			description:      "Should clamp pageSize exceeding max by 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page, pageSize := ValidatePagination(tt.inputPage, tt.inputPageSize)

			if page != tt.expectedPage {
				t.Errorf("%s: page = %d, want %d", tt.description, page, tt.expectedPage)
			}

			if pageSize != tt.expectedPageSize {
				t.Errorf("%s: pageSize = %d, want %d", tt.description, pageSize, tt.expectedPageSize)
			}
		})
	}
}

func TestCalculatePaginationMeta(t *testing.T) {
	tests := []struct {
		name               string
		page               int
		pageSize           int
		total              int64
		expectedTotalPages int
		expectedHasNext    bool
		expectedHasPrev    bool
		description        string
	}{
		{
			name:               "First page with data",
			page:               1,
			pageSize:           20,
			total:              100,
			expectedTotalPages: 5,
			expectedHasNext:    true,
			expectedHasPrev:    false,
			description:        "First page should have next but no previous",
		},
		{
			name:               "Middle page",
			page:               3,
			pageSize:           20,
			total:              100,
			expectedTotalPages: 5,
			expectedHasNext:    true,
			expectedHasPrev:    true,
			description:        "Middle page should have both next and previous",
		},
		{
			name:               "Last page",
			page:               5,
			pageSize:           20,
			total:              100,
			expectedTotalPages: 5,
			expectedHasNext:    false,
			expectedHasPrev:    true,
			description:        "Last page should have previous but no next",
		},
		{
			name:               "Single page with few items",
			page:               1,
			pageSize:           20,
			total:              5,
			expectedTotalPages: 1,
			expectedHasNext:    false,
			expectedHasPrev:    false,
			description:        "Single page should have no navigation",
		},
		{
			name:               "Empty dataset",
			page:               1,
			pageSize:           20,
			total:              0,
			expectedTotalPages: 1,
			expectedHasNext:    false,
			expectedHasPrev:    false,
			description:        "Empty dataset should show 1 total page",
		},
		{
			name:               "Partial last page",
			page:               3,
			pageSize:           20,
			total:              55,
			expectedTotalPages: 3,
			expectedHasNext:    false,
			expectedHasPrev:    true,
			description:        "Should handle partial last page correctly",
		},
		{
			name:               "Exact page boundary",
			page:               2,
			pageSize:           20,
			total:              40,
			expectedTotalPages: 2,
			expectedHasNext:    false,
			expectedHasPrev:    true,
			description:        "Should handle exact page boundary",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := CalculatePaginationMeta(tt.page, tt.pageSize, tt.total)

			if meta.Page != tt.page {
				t.Errorf("%s: Page = %d, want %d", tt.description, meta.Page, tt.page)
			}

			if meta.PageSize != tt.pageSize {
				t.Errorf("%s: PageSize = %d, want %d", tt.description, meta.PageSize, tt.pageSize)
			}

			if meta.Total != tt.total {
				t.Errorf("%s: Total = %d, want %d", tt.description, meta.Total, tt.total)
			}

			if meta.TotalPages != tt.expectedTotalPages {
				t.Errorf("%s: TotalPages = %d, want %d", tt.description, meta.TotalPages, tt.expectedTotalPages)
			}

			if meta.HasNext != tt.expectedHasNext {
				t.Errorf("%s: HasNext = %v, want %v", tt.description, meta.HasNext, tt.expectedHasNext)
			}

			if meta.HasPrev != tt.expectedHasPrev {
				t.Errorf("%s: HasPrev = %v, want %v", tt.description, meta.HasPrev, tt.expectedHasPrev)
			}
		})
	}
}

func TestCalculateOffset(t *testing.T) {
	tests := []struct {
		name           string
		page           int
		pageSize       int
		expectedOffset int
		description    string
	}{
		{
			name:           "First page",
			page:           1,
			pageSize:       20,
			expectedOffset: 0,
			description:    "First page should have offset 0",
		},
		{
			name:           "Second page",
			page:           2,
			pageSize:       20,
			expectedOffset: 20,
			description:    "Second page should skip first page",
		},
		{
			name:           "Page 10 with size 50",
			page:           10,
			pageSize:       50,
			expectedOffset: 450,
			description:    "Should calculate correct offset for large pages",
		},
		{
			name:           "Negative page (security)",
			page:           -1,
			pageSize:       20,
			expectedOffset: 0,
			description:    "Negative page should be normalized to MinPage",
		},
		{
			name:           "Zero page",
			page:           0,
			pageSize:       20,
			expectedOffset: 0,
			description:    "Zero page should be normalized to MinPage",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offset := CalculateOffset(tt.page, tt.pageSize)

			if offset != tt.expectedOffset {
				t.Errorf("%s: offset = %d, want %d", tt.description, offset, tt.expectedOffset)
			}
		})
	}
}

func TestCreateQueryContext(t *testing.T) {
	ctx, cancel := CreateQueryContext()
	defer cancel()

	if ctx == nil {
		t.Error("CreateQueryContext should return non-nil context")
	}

	if cancel == nil {
		t.Error("CreateQueryContext should return non-nil cancel function")
	}

	// Verify context has deadline
	deadline, ok := ctx.Deadline()
	if !ok {
		t.Error("Context should have a deadline")
	}

	// Verify deadline is in the future
	if time.Until(deadline) <= 0 {
		t.Error("Deadline should be in the future")
	}

	// Verify timeout is approximately QueryTimeout
	timeout := time.Until(deadline)
	if timeout < QueryTimeout-time.Second || timeout > QueryTimeout+time.Second {
		t.Errorf("Timeout should be approximately %v, got %v", QueryTimeout, timeout)
	}

	// Verify cancel works
	cancel()
	select {
	case <-ctx.Done():
		// Context was cancelled successfully
		if ctx.Err() != context.Canceled {
			t.Errorf("Expected context.Canceled error, got %v", ctx.Err())
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Context should be cancelled immediately after calling cancel()")
	}
}

func TestPaginationConstants(t *testing.T) {
	// Verify constants have sensible values
	if DefaultPageSize <= 0 {
		t.Error("DefaultPageSize should be positive")
	}

	if MaxPageSize <= DefaultPageSize {
		t.Error("MaxPageSize should be greater than DefaultPageSize")
	}

	if MinPageSize <= 0 {
		t.Error("MinPageSize should be positive")
	}

	if MinPage <= 0 {
		t.Error("MinPage should be positive (1-indexed pagination)")
	}

	if QueryTimeout <= 0 {
		t.Error("QueryTimeout should be positive")
	}

	// Security checks
	if MaxPageSize > 1000 {
		t.Errorf("MaxPageSize (%d) should not exceed 1000 to prevent DoS", MaxPageSize)
	}

	if QueryTimeout > 60*time.Second {
		t.Errorf("QueryTimeout (%v) should not exceed 60 seconds", QueryTimeout)
	}
}

// Benchmark tests to ensure pagination validation is performant
func BenchmarkValidatePagination(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidatePagination(1, 20)
	}
}

func BenchmarkCalculatePaginationMeta(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CalculatePaginationMeta(1, 20, 100)
	}
}

func BenchmarkCalculateOffset(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CalculateOffset(10, 50)
	}
}
