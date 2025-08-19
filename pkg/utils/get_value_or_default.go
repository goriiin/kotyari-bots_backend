package utils

func GetValueOrDefault[T comparable](value, defaultValue T) T {
	var zero T

	if value == zero {
		return defaultValue
	}

	return value
}
