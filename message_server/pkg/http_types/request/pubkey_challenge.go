package request

// Represents a request for a PK challenge.
type PKCRequest struct {
	ID string `json:"id"`
	PK string `json:"pk"`
}

// Represents a signed request for a PK challenge.
type PKCSignedRequest struct {
	PKCRequest `json:",inline" tstype:",extends"`
	Token      string `json:"token"`
	Signature  string `json:"signature"`
}
