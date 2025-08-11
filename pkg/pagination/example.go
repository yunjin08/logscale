package pagination

import (
	"fmt"
)

// ExampleUsage demonstrates how to use the pagination package
func ExampleUsage() {
	// 1. Create pagination from query parameters
	p := NewPaginationFromQuery("2", "25")
	fmt.Printf("Page: %d, Limit: %d, Offset: %d\n", p.Page, p.Limit, p.GetOffset())

	// 2. Set total and get pagination info
	p.SetTotal(100)
	fmt.Printf("Total: %d, Total Pages: %d\n", p.Total, p.TotalPages)
	fmt.Printf("Has Next: %t, Has Prev: %t\n", p.HasNext(), p.HasPrev())

	// 3. Create paginated response
	data := []string{"item1", "item2", "item3"}
	response := CreatePaginatedResponse(data, p)
	fmt.Printf("Response: %+v\n", response)
}

// ExampleMiddlewareUsage shows how to use the middleware
func ExampleMiddlewareUsage() {
	// In your Gin routes:
	// r.GET("/items", pagination.PaginationMiddleware(), yourHandler)
	//
	// In your handler:
	// func yourHandler(c *gin.Context) {
	//     p := pagination.GetPaginationFromContext(c)
	//     // Use p.Page, p.Limit, p.GetOffset() for database queries
	//     // Set total: p.SetTotal(total)
	//     // Return: pagination.CreatePaginatedResponse(data, p)
	// }
}
