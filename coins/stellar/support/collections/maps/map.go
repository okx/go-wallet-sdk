package maps

func Keys[T comparable, U any](m map[T]U) []T {
	keys := make([]T, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func Values[T comparable, U any](m map[T]U) []U {
	values := make([]U, 0, len(m))
	for _, value := range m {
		values = append(values, value)
	}
	return values
}
