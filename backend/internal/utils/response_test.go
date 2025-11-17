package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// setupTestRouter creates a gin router for testing
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// TestSuccessResponse tests the success response function
func TestSuccessResponse(t *testing.T) {
	router := setupTestRouter()
	router.GET("/test", func(c *gin.Context) {
		SuccessResponse(c, http.StatusOK, "Success message", map[string]string{"key": "value"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}

	if response.Message != "Success message" {
		t.Errorf("Expected message 'Success message', got '%s'", response.Message)
	}
}

// TestErrorResponse tests the error response function
func TestErrorResponse(t *testing.T) {
	router := setupTestRouter()
	router.GET("/test", func(c *gin.Context) {
		ErrorResponse(c, http.StatusBadRequest, "ERROR_CODE", "Error message", nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Success {
		t.Error("Expected success to be false")
	}

	if response.Error == nil {
		t.Fatal("Expected error details to be present")
	}

	if response.Error.Code != "ERROR_CODE" {
		t.Errorf("Expected error code 'ERROR_CODE', got '%s'", response.Error.Code)
	}
}

// TestSafeInternalServerErrorResponse tests error sanitization
func TestSafeInternalServerErrorResponse(t *testing.T) {
	tests := []struct {
		name         string
		inputError   error
		expectedMsg  string
		shouldBeInfo bool
	}{
		{
			name:         "GORM record not found",
			inputError:   gorm.ErrRecordNotFound,
			expectedMsg:  ErrMsgNotFound,
			shouldBeInfo: false,
		},
		{
			name:         "Database connection error",
			inputError:   errors.New("dial tcp: connection refused"),
			expectedMsg:  ErrMsgConnectionError,
			shouldBeInfo: false,
		},
		{
			name:         "Generic internal error",
			inputError:   errors.New("unexpected nil pointer dereference at file.go:123"),
			expectedMsg:  ErrMsgInternalServer,
			shouldBeInfo: false,
		},
		{
			name:         "SQL injection attempt",
			inputError:   errors.New("pq: syntax error at or near 'DROP'"),
			expectedMsg:  ErrMsgInternalServer,
			shouldBeInfo: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestRouter()
			router.GET("/test", func(c *gin.Context) {
				SafeInternalServerErrorResponse(c, tt.inputError)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			if w.Code != http.StatusInternalServerError {
				t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
			}

			var response APIResponse
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Success {
				t.Error("Expected success to be false for error response")
			}

			if response.Error == nil {
				t.Fatal("Expected error details to be present")
			}

			// Verify the error message is sanitized (not the raw error)
			if response.Error.Message != tt.expectedMsg {
				t.Errorf("Expected sanitized message '%s', got '%s'", tt.expectedMsg, response.Error.Message)
			}

			// Verify raw error details are NOT exposed
			rawErrMsg := tt.inputError.Error()
			responseBody := w.Body.String()
			if contains(responseBody, rawErrMsg) && rawErrMsg != tt.expectedMsg {
				t.Errorf("Raw error message should not be exposed in response: %s", rawErrMsg)
			}
		})
	}
}

// TestSafeValidationErrorResponse tests validation error sanitization
func TestSafeValidationErrorResponse(t *testing.T) {
	tests := []struct {
		name        string
		inputError  error
		expectedMsg string
	}{
		{
			name:        "Validation error",
			inputError:  errors.New("validation failed: field is required"),
			expectedMsg: ErrMsgValidationError,
		},
		{
			name:        "Invalid format error",
			inputError:  errors.New("invalid email format"),
			expectedMsg: ErrMsgValidationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupTestRouter()
			router.POST("/test", func(c *gin.Context) {
				SafeValidationErrorResponse(c, tt.inputError)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/test", nil)
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
			}

			var response APIResponse
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			if response.Error == nil {
				t.Fatal("Expected error details to be present")
			}

			if response.Error.Message != tt.expectedMsg {
				t.Errorf("Expected message '%s', got '%s'", tt.expectedMsg, response.Error.Message)
			}
		})
	}
}

// TestPaginatedSuccessResponse tests paginated response
func TestPaginatedSuccessResponse(t *testing.T) {
	router := setupTestRouter()
	router.GET("/test", func(c *gin.Context) {
		data := []map[string]string{
			{"id": "1", "name": "Item 1"},
			{"id": "2", "name": "Item 2"},
		}
		pagination := PaginationMeta{
			Page:       1,
			PageSize:   2,
			TotalPages: 5,
			Total:      10,
		}
		PaginatedSuccessResponse(c, http.StatusOK, "Data retrieved", data, &pagination)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response PaginatedResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !response.Success {
		t.Error("Expected success to be true")
	}

	if response.Pagination.Page != 1 {
		t.Errorf("Expected page 1, got %d", response.Pagination.Page)
	}

	if response.Pagination.Total != 10 {
		t.Errorf("Expected total items 10, got %d", response.Pagination.Total)
	}
}

// TestNotFoundResponse tests not found response
func TestNotFoundResponse(t *testing.T) {
	router := setupTestRouter()
	router.GET("/test", func(c *gin.Context) {
		NotFoundResponse(c, "Resource not found")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}

	var response APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Success {
		t.Error("Expected success to be false")
	}

	if response.Error.Code != "NOT_FOUND" {
		t.Errorf("Expected error code 'NOT_FOUND', got '%s'", response.Error.Code)
	}
}

// TestUnauthorizedResponse tests unauthorized response
func TestUnauthorizedResponse(t *testing.T) {
	router := setupTestRouter()
	router.GET("/test", func(c *gin.Context) {
		UnauthorizedResponse(c, "Unauthorized access")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	var response APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Error.Code != "UNAUTHORIZED" {
		t.Errorf("Expected error code 'UNAUTHORIZED', got '%s'", response.Error.Code)
	}
}

// TestCalculatePaginationMeta tests pagination calculation
func TestCalculatePaginationMeta(t *testing.T) {
	tests := []struct {
		name           string
		page           int
		pageSize       int
		totalItems     int64
		expectedPages  int
		expectedItems  int64
	}{
		{
			name:          "Exact pages",
			page:          1,
			pageSize:      10,
			totalItems:    100,
			expectedPages: 10,
			expectedItems: 100,
		},
		{
			name:          "Partial last page",
			page:          2,
			pageSize:      10,
			totalItems:    25,
			expectedPages: 3,
			expectedItems: 25,
		},
		{
			name:          "Single page",
			page:          1,
			pageSize:      20,
			totalItems:    5,
			expectedPages: 1,
			expectedItems: 5,
		},
		{
			name:          "Empty result",
			page:          1,
			pageSize:      10,
			totalItems:    0,
			expectedPages: 0,
			expectedItems: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := CalculatePaginationMeta(tt.page, tt.pageSize, tt.totalItems)

			if meta.Page != tt.page {
				t.Errorf("Expected page %d, got %d", tt.page, meta.Page)
			}

			if meta.PageSize != tt.pageSize {
				t.Errorf("Expected page size %d, got %d", tt.pageSize, meta.PageSize)
			}

			if meta.TotalPages != tt.expectedPages {
				t.Errorf("Expected total pages %d, got %d", tt.expectedPages, meta.TotalPages)
			}

			if meta.Total != tt.expectedItems {
				t.Errorf("Expected total items %d, got %d", tt.expectedItems, meta.Total)
			}
		})
	}
}

// BenchmarkSafeInternalServerErrorResponse benchmarks the safe error response
func BenchmarkSafeInternalServerErrorResponse(b *testing.B) {
	gin.SetMode(gin.TestMode)
	testErr := errors.New("test database error: connection failed")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router := gin.New()
		router.GET("/test", func(c *gin.Context) {
			SafeInternalServerErrorResponse(c, testErr)
		})

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)
	}
}
