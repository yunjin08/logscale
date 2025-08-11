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
	godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize handlers
	logHandler := v1.NewLogHandler(db)

	// Setup Gin router
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r, logHandler)

	log.Println("Starting LogScale API server on :8080")
	r.Run(":8080")
}
