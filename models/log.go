package models

import (
	"encoding/json"
	"time"
)

// Log represents a log entry
type Log struct {
	ID        int64           `json:"id" db:"id"`
	Service   string          `json:"service" db:"service"`
	Level     string          `json:"level" db:"level"`
	Message   string          `json:"message" db:"message"`
	Timestamp time.Time       `json:"timestamp" db:"timestamp"`
	Meta      json.RawMessage `json:"meta" db:"meta"`
}

// LogRequest represents the payload for creating logs
type LogRequest struct {
	Service   string          `json:"service" binding:"required"`
	Level     string          `json:"level" binding:"required"`
	Message   string          `json:"message" binding:"required"`
	Timestamp *time.Time      `json:"timestamp"`
	Meta      json.RawMessage `json:"meta"`
}

// LogQuery represents query parameters for filtering logs
type LogQuery struct {
	Service   string `form:"service"`
	Level     string `form:"level"`
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
}

// LogResponse represents the response for log queries
type LogResponse struct {
	Logs  []Log `json:"logs"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}
