package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	v1 "github.com/yunjin08/logscale/handlers/v1"
	"github.com/yunjin08/logscale/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("warning: failed to load .env file: %v", err)
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Printf("error: DATABASE_URL is not set")
		os.Exit(1)
	}

	db, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Printf("error: failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize handlers
	logHandler := v1.NewLogHandler(db)

	// Setup Gin router
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r, logHandler)

	log.Println("Starting LogScale API server on :8080")
	err = r.Run(":8080")
	if err != nil {
		log.Printf("error: failed to start server: %v", err)
		os.Exit(1)
	}
}
