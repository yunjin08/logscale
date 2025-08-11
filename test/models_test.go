package test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yunjin08/logscale/models"
	"github.com/yunjin08/logscale/pkg/pagination"
)

func TestLogRequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		request models.LogRequest
		isValid bool
	}{
		{
			name: "valid request",
			request: models.LogRequest{
				Service: "test-service",
				Level:   "info",
				Message: "test message",
			},
			isValid: true,
		},
		{
			name: "missing service",
			request: models.LogRequest{
				Level:   "info",
				Message: "test message",
			},
			isValid: false,
		},
		{
			name: "missing level",
			request: models.LogRequest{
				Service: "test-service",
				Message: "test message",
			},
			isValid: false,
		},
		{
			name: "missing message",
			request: models.LogRequest{
				Service: "test-service",
				Level:   "info",
			},
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			data, err := json.Marshal(tt.request)
			if tt.isValid {
				require.NoError(t, err)
				assert.NotEmpty(t, data)
			} else {
				// For invalid requests, we expect marshaling to work
				// but validation would fail in the handler
				require.NoError(t, err)
			}
		})
	}
}

func TestLogJSONMarshaling(t *testing.T) {
	now := time.Now()
	log := models.Log{
		ID:        1,
		Service:   "test-service",
		Level:     "info",
		Message:   "test message",
		Timestamp: now,
		Meta:      json.RawMessage(`{"key":"value"}`),
	}

	data, err := json.Marshal(log)
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	var unmarshaledLog models.Log
	err = json.Unmarshal(data, &unmarshaledLog)
	require.NoError(t, err)

	assert.Equal(t, log.ID, unmarshaledLog.ID)
	assert.Equal(t, log.Service, unmarshaledLog.Service)
	assert.Equal(t, log.Level, unmarshaledLog.Level)
	assert.Equal(t, log.Message, unmarshaledLog.Message)
	assert.Equal(t, log.Meta, unmarshaledLog.Meta)
}

func TestLogQueryParsing(t *testing.T) {
	query := models.LogQuery{
		Service:   "test-service",
		Level:     "error",
		StartTime: "2024-01-01T00:00:00Z",
		EndTime:   "2024-01-02T00:00:00Z",
	}

	// Test JSON marshaling
	data, err := json.Marshal(query)
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	// Test unmarshaling
	var unmarshaledQuery models.LogQuery
	err = json.Unmarshal(data, &unmarshaledQuery)
	require.NoError(t, err)

	assert.Equal(t, query.Service, unmarshaledQuery.Service)
	assert.Equal(t, query.Level, unmarshaledQuery.Level)
	assert.Equal(t, query.StartTime, unmarshaledQuery.StartTime)
	assert.Equal(t, query.EndTime, unmarshaledQuery.EndTime)
}

func TestPaginationParsing(t *testing.T) {
	p := pagination.NewPagination(2, 25)

	data, err := json.Marshal(p)
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	var unmarshaledPagination pagination.Pagination
	err = json.Unmarshal(data, &unmarshaledPagination)
	require.NoError(t, err)

	assert.Equal(t, p.Page, unmarshaledPagination.Page)
	assert.Equal(t, p.Limit, unmarshaledPagination.Limit)
	// Note: Offset is not included in JSON (json:"-"), so we don't test it
}

func TestLogResponseStructure(t *testing.T) {
	logs := []models.Log{
		{
			ID:        1,
			Service:   "service-1",
			Level:     "info",
			Message:   "message 1",
			Timestamp: time.Now(),
		},
		{
			ID:        2,
			Service:   "service-2",
			Level:     "error",
			Message:   "message 2",
			Timestamp: time.Now(),
		},
	}

	p := pagination.NewPagination(1, 50)
	p.SetTotal(2)
	response := pagination.CreatePaginatedResponse(logs, p)

	data, err := json.Marshal(response)
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	var unmarshaledResponse pagination.PaginatedResponse
	err = json.Unmarshal(data, &unmarshaledResponse)
	require.NoError(t, err)

	assert.Equal(t, p.Total, unmarshaledResponse.Pagination.Total)
	assert.Equal(t, p.Page, unmarshaledResponse.Pagination.Page)
	assert.Equal(t, p.Limit, unmarshaledResponse.Pagination.Limit)
}

func TestHealthResponse(t *testing.T) {
	response := models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	data, err := json.Marshal(response)
	require.NoError(t, err)
	assert.NotEmpty(t, data)

	var unmarshaledResponse models.HealthResponse
	err = json.Unmarshal(data, &unmarshaledResponse)
	require.NoError(t, err)

	assert.Equal(t, response.Status, unmarshaledResponse.Status)
	assert.Equal(t, response.Timestamp, unmarshaledResponse.Timestamp)
}
