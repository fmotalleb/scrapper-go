package utils

func MapItems[T any, R any](items []T, mapper func(T) (R, error)) ([]R, error) {
	result := make([]R, len(items))
	for index, value := range items {
		if mappedValue, err := mapper(value); err == nil {
			result[index] = mappedValue
		} else {
			return nil, err
		}
	}
	return result, nil
}
