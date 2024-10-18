package pointer

// Of custom function to get pointer of type
func Of[K comparable](val K) *K {
	return &val
}

// Value return value of pointer. If nil pointer return zero value
func Value[T any](val *T) T {
	if val == nil {
		var zero T
		return zero
	}

	return *val
}
