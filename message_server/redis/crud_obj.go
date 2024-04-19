package redis

// Represents an error returned when querying for multiple Redis keys.
type MultiRedisErr struct {
	error
	indices []int
}

// Gets the cause of a `MultiRedisErr`.
func (er MultiRedisErr) Cause() error {
	return er.error
}

// Gets the list of problematic indices of a `MultiRedisErr`.
func (er MultiRedisErr) Indices() []int {
	return er.indices
}
