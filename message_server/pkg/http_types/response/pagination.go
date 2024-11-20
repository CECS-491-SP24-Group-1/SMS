package response

import (
	"fmt"

	"wraith.me/message_server/pkg/db/qpage"
)

// Represents a paginated collection of objects.
type PaginatedData[T any] struct {
	Data       []T              `json:"data"`
	Pagination qpage.Pagination `json:"pagination"`
}

// Creates a new `PaginatedData` object.
func NewPaginatedData[T any](data []T, pagination qpage.Pagination) PaginatedData[T] {
	return PaginatedData[T]{
		Data:       data,
		Pagination: pagination,
	}
}

// Generates a description string for this page.
func (p PaginatedData[T]) Desc() string {
	typeStr := "<none>"
	if len(p.Data) > 0 {
		typeStr = fmt.Sprintf("%T", p.Data[0])
	}

	page := p.Pagination.CurrentPage
	return fmt.Sprintf("Page %d/%d; containing %d of type %s (%d total)",
		page.Num,
		p.Pagination.TotalPages,
		page.Size,
		typeStr,
		p.Pagination.TotalItems,
	)
}
