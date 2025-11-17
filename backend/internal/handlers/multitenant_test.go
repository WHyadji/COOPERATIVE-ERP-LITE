package handlers_test

import (
	"testing"

	"cooperative-erp-lite/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
)

// TestMultiTenantHelpers tests the multi-tenant helper functions
func TestAmbilIDKoperasiDariContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successfully get koperasi ID from context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		expectedID := uuid.New()
		c.Set("idKoperasi", expectedID)

		id, ok := handlers.AmbilIDKoperasiDariContext(c)
		assert.True(t, ok, "Should successfully get koperasi ID")
		assert.Equal(t, expectedID, id, "Should return correct koperasi ID")
		// When successful, no response is written, so body should be empty
		assert.Empty(t, w.Body.String(), "Should not write any response body")
	})

	t.Run("Fail when koperasi ID not in context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		id, ok := handlers.AmbilIDKoperasiDariContext(c)
		assert.False(t, ok, "Should fail when ID not in context")
		assert.Equal(t, uuid.Nil, id, "Should return nil UUID")
		assert.Equal(t, 500, w.Code, "Should return 500 status code")
	})
}

func TestAmbilIDPenggunaDariContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successfully get pengguna ID from context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		expectedID := uuid.New()
		c.Set("idPengguna", expectedID)

		id, ok := handlers.AmbilIDPenggunaDariContext(c)
		assert.True(t, ok, "Should successfully get pengguna ID")
		assert.Equal(t, expectedID, id, "Should return correct pengguna ID")
		// When successful, no response is written, so body should be empty
		assert.Empty(t, w.Body.String(), "Should not write any response body")
	})

	t.Run("Fail when pengguna ID not in context", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		id, ok := handlers.AmbilIDPenggunaDariContext(c)
		assert.False(t, ok, "Should fail when ID not in context")
		assert.Equal(t, uuid.Nil, id, "Should return nil UUID")
		assert.Equal(t, 401, w.Code, "Should return 401 Unauthorized status code")
	})
}

func TestParseUUIDDariParameter(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successfully parse valid UUID", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		validID := uuid.New()
		c.Params = gin.Params{
			{Key: "id", Value: validID.String()},
		}

		id, ok := handlers.ParseUUIDDariParameter(c, "id")
		assert.True(t, ok, "Should successfully parse UUID")
		assert.Equal(t, validID, id, "Should return correct UUID")
		// When successful, no response is written, so body should be empty
		assert.Empty(t, w.Body.String(), "Should not write any response body")
	})

	t.Run("Fail when UUID is invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Params = gin.Params{
			{Key: "id", Value: "invalid-uuid"},
		}

		id, ok := handlers.ParseUUIDDariParameter(c, "id")
		assert.False(t, ok, "Should fail when UUID is invalid")
		assert.Equal(t, uuid.Nil, id, "Should return nil UUID")
		assert.Equal(t, 400, w.Code, "Should return 400 Bad Request")
	})
}

func TestAmbilParameterPaginasi(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Use default values when no parameters provided", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		params := handlers.AmbilParameterPaginasi(c)
		assert.Equal(t, 1, params.Halaman, "Should use default page 1")
		assert.Equal(t, 20, params.UkuranHalaman, "Should use default page size 20")
	})

	t.Run("Use provided valid parameters", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test?page=5&pageSize=50", nil)

		params := handlers.AmbilParameterPaginasi(c)
		assert.Equal(t, 5, params.Halaman, "Should use provided page")
		assert.Equal(t, 50, params.UkuranHalaman, "Should use provided page size")
	})

	t.Run("Enforce maximum page size limit", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test?page=1&pageSize=200", nil)

		params := handlers.AmbilParameterPaginasi(c)
		assert.Equal(t, 100, params.UkuranHalaman, "Should limit page size to 100")
	})

	t.Run("Handle invalid page number", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test?page=0&pageSize=20", nil)

		params := handlers.AmbilParameterPaginasi(c)
		assert.Equal(t, 1, params.Halaman, "Should reset invalid page to 1")
	})

	t.Run("Handle invalid page size", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test?page=1&pageSize=-5", nil)

		params := handlers.AmbilParameterPaginasi(c)
		assert.Equal(t, 20, params.UkuranHalaman, "Should reset invalid page size to 20")
	})
}
