package try

// Adapter for `go-try` that returns a value if successful or panics otherwise.
func ThrowOnError[T any](val T, err any) T {
	if err != nil {
		panic(err)
	}
	return val
}
