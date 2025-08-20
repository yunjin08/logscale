package stream

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yunjin08/logscale/models"
)

type RedisStreamService struct {
	client     *redis.Client
	streamName string
	config     models.StreamConfig
}

// NewRedisStreamService creates a new Redis Stream service
func NewRedisStreamService(redisURL, streamName string) (*RedisStreamService, error) {
	log.Printf("DEBUG: Initializing Redis Stream service with URL: %s, stream: %s", redisURL, streamName)

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Printf("DEBUG: Failed to parse Redis URL: %v", err)
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("DEBUG: Testing Redis connection...")
	if err = client.Ping(ctx).Err(); err != nil {
		log.Printf("DEBUG: Failed to ping Redis: %v", err)
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}
	log.Printf("DEBUG: Redis connection successful")

	config := models.StreamConfig{
		StreamName:    streamName,
		ConsumerGroup: "logscale-workers",
		ConsumerName:  "worker-1",
	}

	log.Printf("DEBUG: Creating Redis Stream service with config: %+v", config)
	return &RedisStreamService{
		client:     client,
		streamName: streamName,
		config:     config,
	}, nil
}

// PublishLogEvent publishes a log event to Redis Stream
func (s *RedisStreamService) PublishLogEvent(ctx context.Context, logEntry models.Log) error {
	eventMap := map[string]interface{}{
		"id":         fmt.Sprintf("%d", logEntry.ID),
		"service":    logEntry.Service,
		"level":      logEntry.Level,
		"message":    logEntry.Message,
		"timestamp":  logEntry.Timestamp.Format(time.RFC3339),
		"meta":       string(logEntry.Meta),
		"created_at": logEntry.Timestamp.Format(time.RFC3339),
	}

	result := s.client.XAdd(ctx, &redis.XAddArgs{
		Stream: s.streamName,
		Values: eventMap,
	})

	if result.Err() != nil {
		return fmt.Errorf("failed to publish event to stream: %w", result.Err())
	}

	log.Printf("Published event to stream %s: %s", s.streamName, result.Val())
	return nil
}

// CreateConsumerGroup creates a consumer group for the stream
func (s *RedisStreamService) CreateConsumerGroup(ctx context.Context) error {
	// Try to create consumer group, ignore if it already exists
	err := s.client.XGroupCreate(ctx, s.streamName, s.config.ConsumerGroup, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}
	return nil
}

// Close closes the Redis connection
func (s *RedisStreamService) Close() error {
	return s.client.Close()
}

// GetStreamInfo returns information about the stream
func (s *RedisStreamService) GetStreamInfo(ctx context.Context) (*redis.XInfoStream, error) {
	return s.client.XInfoStream(ctx, s.streamName).Result()
}
