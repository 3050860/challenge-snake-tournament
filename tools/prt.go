package tools

func ToPrt[T any](in T) *T {
	return &in
}
