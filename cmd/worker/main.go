package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/yunjin08/logscale/internal/analytics"
	"github.com/yunjin08/logscale/internal/worker"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("warning: failed to load .env file: %v", err)
	}

	// Get environment variables
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Printf("error: DATABASE_URL is not set")
		os.Exit(1)
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Printf("error: REDIS_URL is not set")
		os.Exit(1)
	}

	streamName := os.Getenv("STREAM_NAME")
	if streamName == "" {
		streamName = "logscale:logs" // Default stream name
	}

	// Connect to database
	db, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Printf("error: failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize analytics service
	analyticsSvc := analytics.NewAnalyticsService(db)

	// Initialize worker
	worker, err := worker.NewWorker(redisURL, analyticsSvc, streamName)
	if err != nil {
		log.Printf("error: failed to create worker: %v", err)
		os.Exit(1)
	}
	defer worker.Close()

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start worker in goroutine
	go func() {
		if err := worker.Start(ctx); err != nil {
			log.Printf("error: worker failed: %v", err)
			cancel()
		}
	}()

	log.Printf("Worker started. Press Ctrl+C to stop.")

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutting down worker...")
	cancel()
}
