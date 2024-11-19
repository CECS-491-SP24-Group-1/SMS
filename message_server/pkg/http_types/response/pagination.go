package response

import "wraith.me/message_server/pkg/db/qpage"

// Represents a paginated collection of objects.
type PaginatedData[T any] struct {
	Data       []T              `json:"data"`
	Pagination qpage.Pagination `json:"pagination"`
}
