package analytics

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yunjin08/logscale/models"
)

// AnalyticsService handles metrics aggregation and storage
type AnalyticsService struct {
	db *pgxpool.Pool
}

// NewAnalyticsService creates a new analytics service
func NewAnalyticsService(db *pgxpool.Pool) *AnalyticsService {
	return &AnalyticsService{
		db: db,
	}
}

// UpdateServiceMetrics updates or creates service metrics based on a log event
func (a *AnalyticsService) UpdateServiceMetrics(ctx context.Context, event models.LogEvent) error {
	tx, err := a.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Try to get existing metrics
	var metrics models.ServiceMetrics
	err = tx.QueryRow(ctx, `
		SELECT id, service, total_logs, error_count, warning_count, info_count, debug_count, 
		       error_rate, last_log_time, created_at, updated_at
		FROM service_metrics 
		WHERE service = $1
	`, event.Service).Scan(
		&metrics.ID, &metrics.Service, &metrics.TotalLogs, &metrics.ErrorCount,
		&metrics.WarningCount, &metrics.InfoCount, &metrics.DebugCount,
		&metrics.ErrorRate, &metrics.LastLogTime, &metrics.CreatedAt, &metrics.UpdatedAt,
	)

	if err == sql.ErrNoRows || (err != nil && (err.Error() == "no rows in result set" || strings.Contains(err.Error(), "no rows in result set"))) {
		// Create new metrics record
		metrics = models.ServiceMetrics{
			Service:      event.Service,
			TotalLogs:    0,
			ErrorCount:   0,
			WarningCount: 0,
			InfoCount:    0,
			DebugCount:   0,
			ErrorRate:    0.0,
			LastLogTime:  event.Timestamp,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
	} else if err != nil {
		return fmt.Errorf("failed to query existing metrics: %w", err)
	}

	// Update metrics based on log level
	metrics.TotalLogs++
	metrics.LastLogTime = event.Timestamp
	metrics.UpdatedAt = time.Now()

	switch event.Level {
	case "error":
		metrics.ErrorCount++
	case "warn":
		metrics.WarningCount++
	case "info":
		metrics.InfoCount++
	case "debug":
		metrics.DebugCount++
	}

	// Calculate error rate
	if metrics.TotalLogs > 0 {
		metrics.ErrorRate = float64(metrics.ErrorCount) / float64(metrics.TotalLogs)
	}

	// Insert or update metrics
	if metrics.ID == 0 {
		// Insert new record
		err = tx.QueryRow(ctx, `
			INSERT INTO service_metrics (service, total_logs, error_count, warning_count, info_count, debug_count, error_rate, last_log_time, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id
		`, metrics.Service, metrics.TotalLogs, metrics.ErrorCount, metrics.WarningCount,
			metrics.InfoCount, metrics.DebugCount, metrics.ErrorRate, metrics.LastLogTime,
			metrics.CreatedAt, metrics.UpdatedAt).Scan(&metrics.ID)
	} else {
		// Update existing record
		_, err = tx.Exec(ctx, `
			UPDATE service_metrics 
			SET total_logs = $1, error_count = $2, warning_count = $3, info_count = $4, 
			    debug_count = $5, error_rate = $6, last_log_time = $7, updated_at = $8
			WHERE id = $9
		`, metrics.TotalLogs, metrics.ErrorCount, metrics.WarningCount, metrics.InfoCount,
			metrics.DebugCount, metrics.ErrorRate, metrics.LastLogTime, metrics.UpdatedAt, metrics.ID)
	}

	if err != nil {
		return fmt.Errorf("failed to save metrics: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Updated metrics for service %s: total=%d, errors=%d, error_rate=%.4f",
		metrics.Service, metrics.TotalLogs, metrics.ErrorCount, metrics.ErrorRate)
	return nil
}

// GetServiceMetrics retrieves metrics for a specific service
func (a *AnalyticsService) GetServiceMetrics(ctx context.Context, service string) (*models.ServiceMetrics, error) {
	var metrics models.ServiceMetrics
	err := a.db.QueryRow(ctx, `
		SELECT id, service, total_logs, error_count, warning_count, info_count, debug_count, 
		       error_rate, last_log_time, created_at, updated_at
		FROM service_metrics 
		WHERE service = $1
	`, service).Scan(
		&metrics.ID, &metrics.Service, &metrics.TotalLogs, &metrics.ErrorCount,
		&metrics.WarningCount, &metrics.InfoCount, &metrics.DebugCount,
		&metrics.ErrorRate, &metrics.LastLogTime, &metrics.CreatedAt, &metrics.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no metrics found for service: %s", service)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query metrics: %w", err)
	}

	return &metrics, nil
}

// GetAllServiceMetrics retrieves metrics for all services
func (a *AnalyticsService) GetAllServiceMetrics(ctx context.Context) ([]models.ServiceMetrics, error) {
	rows, err := a.db.Query(ctx, `
		SELECT id, service, total_logs, error_count, warning_count, info_count, debug_count, 
		       error_rate, last_log_time, created_at, updated_at
		FROM service_metrics 
		ORDER BY service
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query metrics: %w", err)
	}
	defer rows.Close()

	var metrics []models.ServiceMetrics
	for rows.Next() {
		var m models.ServiceMetrics
		err := rows.Scan(
			&m.ID, &m.Service, &m.TotalLogs, &m.ErrorCount,
			&m.WarningCount, &m.InfoCount, &m.DebugCount,
			&m.ErrorRate, &m.LastLogTime, &m.CreatedAt, &m.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan metrics: %w", err)
		}
		metrics = append(metrics, m)
	}

	return metrics, nil
}
