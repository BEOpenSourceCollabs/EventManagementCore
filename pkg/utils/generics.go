package utils

// AsPtr returns a pointer to the type given as argument t.
func AsPtr[T any](t T) *T {
	return &t
}
