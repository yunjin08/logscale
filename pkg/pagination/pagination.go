package pagination

import (
	"strconv"
)

// Pagination represents pagination parameters and results
type Pagination struct {
	Page       int   `json:"page" form:"page"`
	Limit      int   `json:"limit" form:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	Offset     int   `json:"-"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

// DefaultPagination returns default pagination values
func DefaultPagination() Pagination {
	return Pagination{
		Page:  1,
		Limit: 50,
	}
}

// NewPagination creates a new pagination instance with validation
func NewPagination(page, limit int) Pagination {
	p := Pagination{
		Page:  page,
		Limit: limit,
	}
	p.Validate()
	return p
}

// NewPaginationFromQuery creates pagination from query parameters
func NewPaginationFromQuery(pageStr, limitStr string) Pagination {
	p := DefaultPagination()

	if pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			p.Page = page
		}
	}

	if limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			p.Limit = limit
		}
	}

	p.Validate()
	return p
}

// Validate validates and corrects pagination parameters
func (p *Pagination) Validate() {
	// Validate page
	if p.Page < 1 {
		p.Page = 1
	}

	// Validate limit
	if p.Limit < 1 {
		p.Limit = 50
	}
	if p.Limit > 1000 {
		p.Limit = 1000
	}

	// Calculate offset
	p.Offset = (p.Page - 1) * p.Limit
}

// SetTotal sets the total count and calculates total pages
func (p *Pagination) SetTotal(total int64) {
	p.Total = total
	if p.Limit > 0 {
		p.TotalPages = int((total + int64(p.Limit) - 1) / int64(p.Limit))
	}
}

// HasNext returns true if there's a next page
func (p *Pagination) HasNext() bool {
	return p.Page < p.TotalPages
}

// HasPrev returns true if there's a previous page
func (p *Pagination) HasPrev() bool {
	return p.Page > 1
}

// NextPage returns the next page number
func (p *Pagination) NextPage() int {
	if p.HasNext() {
		return p.Page + 1
	}
	return p.Page
}

// PrevPage returns the previous page number
func (p *Pagination) PrevPage() int {
	if p.HasPrev() {
		return p.Page - 1
	}
	return p.Page
}

// GetOffset returns the calculated offset
func (p *Pagination) GetOffset() int {
	return p.Offset
}

// GetLimit returns the limit
func (p *Pagination) GetLimit() int {
	return p.Limit
}

// ToMap converts pagination to a map for query parameters
func (p *Pagination) ToMap() map[string]string {
	return map[string]string{
		"page":  strconv.Itoa(p.Page),
		"limit": strconv.Itoa(p.Limit),
	}
}
