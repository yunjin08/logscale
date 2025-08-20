package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yunjin08/logscale/internal/analytics"
	"github.com/yunjin08/logscale/models"
)

// Worker processes events from Redis Streams
type Worker struct {
	redisClient   *redis.Client
	analyticsSvc  *analytics.Service
	streamName    string
	consumerGroup string
	consumerName  string
	maxRetries    int
	retryDelay    time.Duration
}

// NewWorker creates a new worker instance
func NewWorker(redisURL string, analyticsSvc *analytics.Service, streamName string) (*Worker, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Worker{
		redisClient:   client,
		analyticsSvc:  analyticsSvc,
		streamName:    streamName,
		consumerGroup: "logscale-workers",
		consumerName:  "worker-1",
		maxRetries:    3,
		retryDelay:    5 * time.Second,
	}, nil
}

// Start begins processing events from the stream
func (w *Worker) Start(ctx context.Context) error {
	// Create consumer group if it doesn't exist
	err := w.createConsumerGroup(ctx)
	if err != nil {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}

	log.Printf("Worker started. Listening to stream: %s, group: %s, consumer: %s",
		w.streamName, w.consumerGroup, w.consumerName)

	// Process events in a loop
	for {
		select {
		case <-ctx.Done():
			log.Println("Worker stopped by context cancellation")
			return nil
		default:
			err := w.processEvents(ctx)
			if err != nil {
				log.Printf("Error processing events: %v", err)
				time.Sleep(w.retryDelay)
			}
		}
	}
}

// createConsumerGroup creates the consumer group for the stream
func (w *Worker) createConsumerGroup(ctx context.Context) error {
	// Try to create consumer group with MKSTREAM option to create stream if it doesn't exist
	err := w.redisClient.XGroupCreateMkStream(ctx, w.streamName, w.consumerGroup, "0").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}
	return nil
}

// processEvents reads and processes events from the stream
func (w *Worker) processEvents(ctx context.Context) error {
	// Read events from the stream
	streams, err := w.redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    w.consumerGroup,
		Consumer: w.consumerName,
		Streams:  []string{w.streamName, ">"},
		Count:    10,
		Block:    1 * time.Second,
	}).Result()

	if err != nil && err != redis.Nil {
		return fmt.Errorf("failed to read from stream: %w", err)
	}

	if len(streams) == 0 {
		return nil // No events to process
	}

	// Process each stream
	for _, stream := range streams {
		for _, message := range stream.Messages {

			err := w.processEvent(ctx, message)
			if err != nil {
				log.Printf("Failed to process event %s: %v", message.ID, err)
				// Acknowledge the message to prevent reprocessing
				if ackErr := w.acknowledgeMessage(ctx, message.ID); ackErr != nil {
					log.Printf("Failed to acknowledge failed message %s: %v", message.ID, ackErr)
				}
				continue
			}

			// Acknowledge successful processing
			err = w.acknowledgeMessage(ctx, message.ID)
			if err != nil {
				log.Printf("Failed to acknowledge message %s: %v", message.ID, err)
			}
		}
	}

	return nil
}

// processEvent processes a single event
func (w *Worker) processEvent(ctx context.Context, message redis.XMessage) error {
	// Parse the event from message values
	event, err := w.parseEvent(message.Values)
	if err != nil {
		return fmt.Errorf("failed to parse event: %w", err)
	}

	err = w.analyticsSvc.UpdateServiceMetrics(ctx, *event)
	if err != nil {
		return fmt.Errorf("failed to update analytics: %w", err)
	}

	log.Printf("Processed event: service=%s, level=%s, id=%s",
		event.Service, event.Level, event.ID)
	return nil
}

// parseEvent converts Redis message values to LogEvent
func (w *Worker) parseEvent(values map[string]interface{}) (*models.LogEvent, error) {
	// Extract values from the message
	id, ok := values["id"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid id field")
	}

	service, ok := values["service"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid service field")
	}

	level, ok := values["level"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid level field")
	}

	message, ok := values["message"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid message field")
	}

	timestampStr, ok := values["timestamp"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid timestamp field")
	}

	timestamp, err := time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		return nil, fmt.Errorf("invalid timestamp format: %w", err)
	}

	metaStr, ok := values["meta"].(string)
	if !ok {
		metaStr = "{}"
	}

	meta := json.RawMessage(metaStr)

	createdAtStr, ok := values["created_at"].(string)
	if !ok {
		createdAtStr = timestampStr
	}

	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		createdAt = timestamp
	}

	return &models.LogEvent{
		ID:        id,
		Service:   service,
		Level:     level,
		Message:   message,
		Timestamp: timestamp,
		Meta:      meta,
		CreatedAt: createdAt,
	}, nil
}

// acknowledgeMessage acknowledges a processed message
func (w *Worker) acknowledgeMessage(ctx context.Context, messageID string) error {
	return w.redisClient.XAck(ctx, w.streamName, w.consumerGroup, messageID).Err()
}

// Close closes the Redis connection
func (w *Worker) Close() error {
	return w.redisClient.Close()
}
