-- Create logs table
CREATE TABLE logs (
    id BIGSERIAL PRIMARY KEY,
    service VARCHAR(255) NOT NULL,
    level VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    meta JSONB,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Create indexes for better query performance
CREATE INDEX idx_logs_service ON logs(service);
CREATE INDEX idx_logs_level ON logs(level);
CREATE INDEX idx_logs_timestamp ON logs(timestamp);
CREATE INDEX idx_logs_service_level ON logs(service, level);
CREATE INDEX idx_logs_timestamp_desc ON logs(timestamp DESC);

-- Add comments
COMMENT ON TABLE logs IS 'Stores application logs with metadata';
COMMENT ON COLUMN logs.id IS 'Unique identifier for the log entry';
COMMENT ON COLUMN logs.service IS 'Name of the service that generated the log';
COMMENT ON COLUMN logs.level IS 'Log level (debug, info, warn, error, fatal)';
COMMENT ON COLUMN logs.message IS 'The log message content';
COMMENT ON COLUMN logs.timestamp IS 'When the log was generated';
COMMENT ON COLUMN logs.meta IS 'Additional metadata as JSON';
COMMENT ON COLUMN logs.created_at IS 'When the log was stored in the database';
