package util

func SetIfNil[T int | int32 | int64 | string | float32 | float64](value *T, set T) T {
	if value != nil {
		return *value
	}
	return set
}

func SetIf0[T int | int32 | int64 | float32 | float64](value T, set T) T {
	if value != T(0) {
		return value
	}
	return set
}

func SetIfEmpty(value string, set string) string {
	if value != "" {
		return value
	}
	return set
}
