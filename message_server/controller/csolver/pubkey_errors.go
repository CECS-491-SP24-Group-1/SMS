package csolver

import "errors"

var (
	//The error to throw when the pubkey challenge request fails to be parsed.
	PKReqParseError = errors.New("failed to parse pubkey challenge request with error: %s")
)
