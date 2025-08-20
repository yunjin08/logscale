package v1

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yunjin08/logscale/helpers"
	"github.com/yunjin08/logscale/internal/stream"
	"github.com/yunjin08/logscale/models"
	"github.com/yunjin08/logscale/pkg/pagination"
)

type LogHandler struct {
	db        *pgxpool.Pool
	helper    *helpers.LogHelper
	streamSvc *stream.RedisStreamService
}

func NewLogHandler(db *pgxpool.Pool, streamSvc *stream.RedisStreamService) *LogHandler {
	return &LogHandler{
		db:        db,
		helper:    helpers.NewLogHelper(db),
		streamSvc: streamSvc,
	}
}

// CreateLog handles
// POST /v1/logs - accepts batch or single log payload
func (h *LogHandler) CreateLog(c *gin.Context) {
	var request struct {
		Logs []models.LogRequest `json:"logs"`
		Log  *models.LogRequest  `json:"log"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// Handle single log
	if request.Log != nil {
		log, err := h.helper.CreateSingleLog(c.Request.Context(), *request.Log)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Publish to Redis Stream (async)
		go func() {
			if h.streamSvc != nil {
				if err := h.streamSvc.PublishLogEvent(context.Background(), *log); err != nil {
					// Log error but don't fail the request
					_ = c.Error(err) // Ignore error return from c.Error
				}
			}
		}()

		c.JSON(http.StatusCreated, log)
		return
	}

	// Handle batch logs
	if len(request.Logs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No logs provided"})
		return
	}

	logs, err := h.helper.CreateBatchLogs(c.Request.Context(), request.Logs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Publish batch to Redis Stream (async)
	go func() {
		if h.streamSvc != nil {
			for _, log := range logs {
				if err := h.streamSvc.PublishLogEvent(context.Background(), log); err != nil {
					// Log error but don't fail the request
					_ = c.Error(err) // Ignore error return from c.Error
				}
			}
		}
	}()

	c.JSON(http.StatusCreated, gin.H{"logs": logs, "count": len(logs)})
}

// GetLogs handles
// GET /v1/logs - query by service/level/time (paginated)
func (h *LogHandler) GetLogs(c *gin.Context) {
	var query models.LogQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	// Get pagination from context (set by middleware)
	p := pagination.GetPaginationFromContext(c)

	logs, total, err := h.helper.QueryLogs(c.Request.Context(), query, p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set total and create paginated response
	p.SetTotal(total)
	response := pagination.CreatePaginatedResponse(logs, p)

	c.JSON(http.StatusOK, response)
}

// Health handles
// GET /health - readiness/liveness
func (h *LogHandler) Health(c *gin.Context) {
	// Check database connectivity
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err := h.helper.PingDatabase(ctx)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, models.HealthResponse{
			Status:    "unhealthy",
			Timestamp: time.Now().Format(time.RFC3339),
		})
		return
	}

	c.JSON(http.StatusOK, models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().Format(time.RFC3339),
	})
}
