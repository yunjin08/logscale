package test

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yunjin08/logscale/models"
	"github.com/yunjin08/logscale/pkg/pagination"
)

// SetupTestRouter creates a test router with mock handlers
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Mock handler for testing
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, models.HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now().Format(time.RFC3339),
		})
	})

	r.POST("/v1/logs", func(c *gin.Context) {
		var request struct {
			Logs []models.LogRequest `json:"logs"`
			Log  *models.LogRequest  `json:"log"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		if request.Log != nil {
			// Mock single log response
			log := models.Log{
				ID:        1,
				Service:   request.Log.Service,
				Level:     request.Log.Level,
				Message:   request.Log.Message,
				Timestamp: time.Now(),
				Meta:      request.Log.Meta,
			}
			c.JSON(http.StatusCreated, log)
			return
		}

		if len(request.Logs) > 0 {
			// Mock batch response
			logs := make([]models.Log, len(request.Logs))
			for i, req := range request.Logs {
				logs[i] = models.Log{
					ID:        int64(i + 1),
					Service:   req.Service,
					Level:     req.Level,
					Message:   req.Message,
					Timestamp: time.Now(),
					Meta:      req.Meta,
				}
			}
			c.JSON(http.StatusCreated, gin.H{"logs": logs, "count": len(logs)})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": "No logs provided"})
	})

	r.GET("/v1/logs", func(c *gin.Context) {
		// Mock query response with pagination
		logs := []models.Log{
			{
				ID:        1,
				Service:   "test-service",
				Level:     "info",
				Message:   "Test log message",
				Timestamp: time.Now(),
				Meta:      json.RawMessage(`{"test": "data"}`),
			},
		}

		// Create paginated response
		p := pagination.NewPagination(1, 50)
		p.SetTotal(1)
		response := pagination.CreatePaginatedResponse(logs, p)

		c.JSON(http.StatusOK, response)
	})

	return r
}

// CreateTestLogRequest creates a test log request
func CreateTestLogRequest(service, level, message string) models.LogRequest {
	return models.LogRequest{
		Service: service,
		Level:   level,
		Message: message,
		Meta:    json.RawMessage(`{"test":"data"}`),
	}
}

// CreateTestLog creates a test log
func CreateTestLog(id int64, service, level, message string) models.Log {
	return models.Log{
		ID:        id,
		Service:   service,
		Level:     level,
		Message:   message,
		Timestamp: time.Now(),
		Meta:      json.RawMessage(`{"test":"data"}`),
	}
}
