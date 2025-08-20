package models

import (
	"encoding/json"
	"time"
)

// LogEvent represents an event published to Redis Stream
type LogEvent struct {
	ID        string          `json:"id"`
	Service   string          `json:"service"`
	Level     string          `json:"level"`
	Message   string          `json:"message"`
	Timestamp time.Time       `json:"timestamp"`
	Meta      json.RawMessage `json:"meta"`
	CreatedAt time.Time       `json:"created_at"`
}

// StreamConfig holds Redis Stream configuration
type StreamConfig struct {
	StreamName    string `json:"stream_name"`
	ConsumerGroup string `json:"consumer_group"`
	ConsumerName  string `json:"consumer_name"`
}

// ServiceMetrics represents aggregated analytics
type ServiceMetrics struct {
	ID           int64     `json:"id" db:"id"`
	Service      string    `json:"service" db:"service"`
	TotalLogs    int64     `json:"total_logs" db:"total_logs"`
	ErrorCount   int64     `json:"error_count" db:"error_count"`
	WarningCount int64     `json:"warning_count" db:"warning_count"`
	InfoCount    int64     `json:"info_count" db:"info_count"`
	DebugCount   int64     `json:"debug_count" db:"debug_count"`
	ErrorRate    float64   `json:"error_rate" db:"error_rate"`
	LastLogTime  time.Time `json:"last_log_time" db:"last_log_time"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// DeadLetterEvent represents failed events
type DeadLetterEvent struct {
	ID         string          `json:"id"`
	OriginalID string          `json:"original_id"`
	Event      json.RawMessage `json:"event"`
	Error      string          `json:"error"`
	RetryCount int             `json:"retry_count"`
	FailedAt   time.Time       `json:"failed_at"`
	StreamName string          `json:"stream_name"`
}
