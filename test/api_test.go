package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yunjin08/logscale/models"
	"github.com/yunjin08/logscale/pkg/pagination"
)

func TestHealthCheck(t *testing.T) {
	router := SetupTestRouter()

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

func TestCreateSingleLog(t *testing.T) {
	router := SetupTestRouter()

	logRequest := CreateTestLogRequest("test-service", "info", "Test log message")

	requestBody := gin.H{
		"log": logRequest,
	}

	jsonData, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/v1/logs", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Log
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, logRequest.Service, response.Service)
	assert.Equal(t, logRequest.Level, response.Level)
	assert.Equal(t, logRequest.Message, response.Message)
	assert.Equal(t, logRequest.Meta, response.Meta)
}

func TestCreateBatchLogs(t *testing.T) {
	router := SetupTestRouter()

	logRequests := []models.LogRequest{
		CreateTestLogRequest("service-1", "info", "First log message"),
		CreateTestLogRequest("service-2", "error", "Second log message"),
	}

	requestBody := gin.H{
		"logs": logRequests,
	}

	jsonData, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/v1/logs", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response struct {
		Logs  []models.Log `json:"logs"`
		Count int          `json:"count"`
	}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, 2, response.Count)
	assert.Len(t, response.Logs, 2)
	assert.Equal(t, logRequests[0].Service, response.Logs[0].Service)
	assert.Equal(t, logRequests[1].Service, response.Logs[1].Service)
}

func TestCreateLogsInvalidRequest(t *testing.T) {
	router := SetupTestRouter()

	// Test with empty request
	req, err := http.NewRequest("POST", "/v1/logs", bytes.NewBuffer([]byte("{}")))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetLogs(t *testing.T) {
	router := SetupTestRouter()

	req, err := http.NewRequest("GET", "/v1/logs", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response pagination.PaginatedResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, int64(1), response.Pagination.Total)
	assert.Equal(t, 1, response.Pagination.Page)
	assert.Equal(t, 50, response.Pagination.Limit)

	// Check data ser
	var logs []models.Log
	logsData, _ := json.Marshal(response.Data)
	err = json.Unmarshal(logsData, &logs)
	require.NoError(t, err)
	assert.Len(t, logs, 1)
	assert.Equal(t, "test-service", logs[0].Service)
}

func TestGetLogsWithQueryParams(t *testing.T) {
	router := SetupTestRouter()

	req, err := http.NewRequest("GET", "/v1/logs?service=test-service&level=info&page=1&limit=20", nil)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response pagination.PaginatedResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, int64(1), response.Pagination.Total)
	assert.Equal(t, 1, response.Pagination.Page)
	assert.Equal(t, 50, response.Pagination.Limit) // Mock returns default limit
}
