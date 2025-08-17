-- Create service_metrics table for aggregated analytics
CREATE TABLE service_metrics (
    id BIGSERIAL PRIMARY KEY,
    service VARCHAR(255) NOT NULL,
    total_logs BIGINT NOT NULL DEFAULT 0,
    error_count BIGINT NOT NULL DEFAULT 0,
    warning_count BIGINT NOT NULL DEFAULT 0,
    info_count BIGINT NOT NULL DEFAULT 0,
    debug_count BIGINT NOT NULL DEFAULT 0,
    error_rate DECIMAL(5,4) NOT NULL DEFAULT 0.0000,
    last_log_time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(service)
);

-- Create indexes for better query performance
CREATE INDEX idx_service_metrics_service ON service_metrics(service);
CREATE INDEX idx_service_metrics_error_rate ON service_metrics(error_rate);
CREATE INDEX idx_service_metrics_last_log_time ON service_metrics(last_log_time);

-- Create dead_letter_events table for failed events
CREATE TABLE dead_letter_events (
    id BIGSERIAL PRIMARY KEY,
    original_id VARCHAR(255) NOT NULL,
    event JSONB NOT NULL,
    error TEXT NOT NULL,
    retry_count INTEGER NOT NULL DEFAULT 0,
    failed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    stream_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for dead letter events
CREATE INDEX idx_dead_letter_events_stream_name ON dead_letter_events(stream_name);
CREATE INDEX idx_dead_letter_events_failed_at ON dead_letter_events(failed_at);
CREATE INDEX idx_dead_letter_events_retry_count ON dead_letter_events(retry_count);

-- Add comments
COMMENT ON TABLE service_metrics IS 'Aggregated metrics per service';
COMMENT ON COLUMN service_metrics.service IS 'Name of the service';
COMMENT ON COLUMN service_metrics.total_logs IS 'Total number of logs for this service';
COMMENT ON COLUMN service_metrics.error_count IS 'Number of error level logs';
COMMENT ON COLUMN service_metrics.warning_count IS 'Number of warning level logs';
COMMENT ON COLUMN service_metrics.info_count IS 'Number of info level logs';
COMMENT ON COLUMN service_metrics.debug_count IS 'Number of debug level logs';
COMMENT ON COLUMN service_metrics.error_rate IS 'Error rate as decimal (0.0000 to 1.0000)';
COMMENT ON COLUMN service_metrics.last_log_time IS 'Timestamp of the most recent log';

COMMENT ON TABLE dead_letter_events IS 'Failed events that could not be processed';
COMMENT ON COLUMN dead_letter_events.original_id IS 'Original event ID from the stream';
COMMENT ON COLUMN dead_letter_events.event IS 'The original event data as JSON';
COMMENT ON COLUMN dead_letter_events.error IS 'Error message explaining why processing failed';
COMMENT ON COLUMN dead_letter_events.retry_count IS 'Number of times this event was retried';
COMMENT ON COLUMN dead_letter_events.stream_name IS 'Name of the stream this event came from'; 