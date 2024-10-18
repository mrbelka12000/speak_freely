package pointer

func Of[K comparable](val K) *K {
	return &val
}

func Value[T any](val *T) T {
	if val == nil {
		var zero T
		return zero
	}

	return *val
}
