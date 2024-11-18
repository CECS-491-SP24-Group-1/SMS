package qpage

// Represents all of the pages in the paginated query.
type Pagination struct {
	CurrentPage Page  `json:"current_page"`
	PerPage     int   `json:"per_page"`
	TotalPages  int64 `json:"total_pages"`
	TotalItems  int64 `json:"total_items"`
}

// Represents a single page in a pagination query.
type Page struct {
	Num      int    `json:"num"`
	Size     int    `json:"size"`
	IsLast   bool   `json:"is_last"`
	IsEmpty  bool   `json:"is_empty"`
	FirstIdx int    `json:"first_idx"`
	LastIdx  int    `json:"last_idx"`
	FirstID  string `json:"first_id,omitempty"` // ID of the first item on current page
	LastID   string `json:"last_id,omitempty"`  // ID of the last item on current page
}
