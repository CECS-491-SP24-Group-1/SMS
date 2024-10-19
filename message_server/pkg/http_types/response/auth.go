package response

// Represents an info object passed when logging in or refreshing tokens.
type Auth struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// Represents a login token issued during a login request.
type LoginReq struct {
	Token string `json:"token"`
}
