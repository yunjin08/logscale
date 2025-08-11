package helpers

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yunjin08/logscale/models"
	"github.com/yunjin08/logscale/pkg/pagination"
)

// LogHelper contains database operations for logs
type LogHelper struct {
	db *pgxpool.Pool
}

// NewLogHelper creates a new LogHelper instance
func NewLogHelper(db *pgxpool.Pool) *LogHelper {
	return &LogHelper{db: db}
}

// CreateSingleLog creates a single log entry in the database
func (h *LogHelper) CreateSingleLog(ctx context.Context, req models.LogRequest) (*models.Log, error) {
	timestamp := time.Now()
	if req.Timestamp != nil {
		timestamp = *req.Timestamp
	}

	query := `
		INSERT INTO logs (service, level, message, timestamp, meta)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, service, level, message, timestamp, meta
	`

	var log models.Log
	err := h.db.QueryRow(ctx, query,
		req.Service,
		req.Level,
		req.Message,
		timestamp,
		req.Meta,
	).Scan(&log.ID, &log.Service, &log.Level, &log.Message, &log.Timestamp, &log.Meta)

	if err != nil {
		return nil, fmt.Errorf("failed to create log: %w", err)
	}

	return &log, nil
}

// CreateBatchLogs creates multiple log entries in a single transaction
func (h *LogHelper) CreateBatchLogs(ctx context.Context, requests []models.LogRequest) ([]models.Log, error) {
	tx, err := h.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO logs (service, level, message, timestamp, meta)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, service, level, message, timestamp, meta
	`

	var logs []models.Log
	for _, req := range requests {
		timestamp := time.Now()
		if req.Timestamp != nil {
			timestamp = *req.Timestamp
		}

		var log models.Log
		err := tx.QueryRow(ctx, query,
			req.Service,
			req.Level,
			req.Message,
			timestamp,
			req.Meta,
		).Scan(&log.ID, &log.Service, &log.Level, &log.Message, &log.Timestamp, &log.Meta)

		if err != nil {
			return nil, fmt.Errorf("failed to create log in batch: %w", err)
		}

		logs = append(logs, log)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return logs, nil
}

// QueryLogs retrieves logs with filtering and pagination
func (h *LogHelper) QueryLogs(ctx context.Context, query models.LogQuery, pagination pagination.Pagination) ([]models.Log, int64, error) {
	// Build WHERE clause
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argCount := 1

	if query.Service != "" {
		whereClause += fmt.Sprintf(" AND service = $%d", argCount)
		args = append(args, query.Service)
		argCount++
	}

	if query.Level != "" {
		whereClause += fmt.Sprintf(" AND level = $%d", argCount)
		args = append(args, query.Level)
		argCount++
	}

	if query.StartTime != "" {
		whereClause += fmt.Sprintf(" AND timestamp >= $%d", argCount)
		args = append(args, query.StartTime)
		argCount++
	}

	if query.EndTime != "" {
		whereClause += fmt.Sprintf(" AND timestamp <= $%d", argCount)
		args = append(args, query.EndTime)
		argCount++
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM logs %s", whereClause)
	var total int64
	err := h.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count logs: %w", err)
	}

	// Get paginated results using pagination package
	args = append(args, pagination.GetLimit(), pagination.GetOffset())

	dataQuery := fmt.Sprintf(`
		SELECT id, service, level, message, timestamp, meta
		FROM logs %s
		ORDER BY timestamp DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argCount, argCount+1)

	rows, err := h.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query logs: %w", err)
	}
	defer rows.Close()

	var logs []models.Log
	for rows.Next() {
		var log models.Log
		err := rows.Scan(&log.ID, &log.Service, &log.Level, &log.Message, &log.Timestamp, &log.Meta)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan log: %w", err)
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating logs: %w", err)
	}

	return logs, total, nil
}

// PingDatabase checks database connectivity
func (h *LogHelper) PingDatabase(ctx context.Context) error {
	return h.db.Ping(ctx)
}
