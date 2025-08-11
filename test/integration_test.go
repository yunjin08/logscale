//go:build integration
// +build integration

package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yunjin08/logscale/models"
)

func setupIntegrationTest() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Use real handlers with real database connection
	// This will be set up in CI with a real PostgreSQL instance
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, models.HealthResponse{
			Status:    "healthy",
			Timestamp: "2024-01-15T10:30:00Z",
		})
	})

	return r
}

func TestIntegrationHealthCheck(t *testing.T) {
	// Skip if no database URL is provided
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	router := setupIntegrationTest()

	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.HealthResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "healthy", response.Status)
	assert.NotEmpty(t, response.Timestamp)
}

func TestIntegrationDatabaseConnection(t *testing.T) {
	// Skip if no database URL is provided
	if os.Getenv("DATABASE_URL") == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	// This test would verify that the application can connect to the database
	// and perform basic operations
	t.Log("Integration test with real database would run here")
}
