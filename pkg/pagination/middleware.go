package pagination

import (
	"github.com/gin-gonic/gin"
)

// PaginationMiddleware extracts pagination parameters from query string
func PaginationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		pageStr := c.Query("page")
		limitStr := c.Query("limit")

		pagination := NewPaginationFromQuery(pageStr, limitStr)

		// Store pagination in context for handlers to use
		c.Set("pagination", pagination)

		c.Next()
	}
}

// GetPaginationFromContext retrieves pagination from Gin context
func GetPaginationFromContext(c *gin.Context) Pagination {
	if pagination, exists := c.Get("pagination"); exists {
		if p, ok := pagination.(Pagination); ok {
			return p
		}
	}
	return DefaultPagination()
}

// CreatePaginatedResponse creates a standardized paginated response
func CreatePaginatedResponse(data interface{}, pagination Pagination) PaginatedResponse {
	return PaginatedResponse{
		Data:       data,
		Pagination: pagination,
	}
}
