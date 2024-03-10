package obj

// Represents an object that's returned to a client when an error occurs.
type ErrorResp struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Error  string `json:"error"`
}

// Represents an object that's returned to a client when multiple errors occur.
type MultiErrorResp struct {
	Status string   `json:"status"`
	Code   int      `json:"code"`
	Errors []string `json:"errors"`
}
