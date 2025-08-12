package routes

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/yunjin08/logscale/handlers/v1"
	"github.com/yunjin08/logscale/pkg/pagination"
)

// SetupRoutes configures all the API routes
func SetupRoutes(r *gin.Engine, logHandler *v1.LogHandler) {
	// Health check endpoint
	r.GET("/health", logHandler.Health)

	// API v1 routes
	v1 := r.Group("/v1")
	{
		// Logs endpoints
		logs := v1.Group("/logs")
		{
			logs.POST("", logHandler.CreateLog)                       // POST /v1/logs
			logs.GET("", pagination.Middleware(), logHandler.GetLogs) // GET /v1/logs with pagination
		}
	}

	// Root endpoint for basic info
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "LogScale API",
			"version": "v1",
			"endpoints": gin.H{
				"health": "/health",
				"logs":   "/v1/logs",
			},
		})
	})
}
