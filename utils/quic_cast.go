package utils

func CastItems[T any](items []any) ([]T, bool) {
	result := make([]T, len(items))
	for index, value := range items {
		if value, ok := value.(T); ok {
			result[index] = value
		} else {
			return nil, false
		}
	}
	return result, true
}
