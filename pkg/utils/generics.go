package utils

func AsPtr[T any](t T) *T {
	return &t
}
