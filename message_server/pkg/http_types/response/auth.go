package response

import "time"

// Represents an info object passed when logging in or refreshing tokens.
type Auth struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

// Represents a login token issued during a login request.
type LoginReq struct {
	Token string `json:"token"`
}

// Defines the structure of a session had via querying for a user's refresh tokens.
type Session struct {
	ID        string    `json:"id"`
	IsCurrent bool      `json:"is_current"`
	Created   time.Time `json:"created"`
	Expires   time.Time `json:"expires"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"string"`
}

// Represents a single access token had via getting the current session info.
type AccessSession struct {
	ID      string    `json:"id"`
	Created time.Time `json:"created"`
	Expires time.Time `json:"expires"`
	Parent  Session   `json:"parent"`
}

// Represents a list of sessions, keyed by the session's ID.
type SessionsList map[string]Session
