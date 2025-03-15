package utils

import "fmt"

func SingleOrMulti[T any](items map[string]any, key string) []T {
	result := make([]T, 0)
	if value, ok := items[key]; ok {
		if values, err := SliceOf[T](value); err == nil {
			result = append(result, values...)
		}
	}
	keys := fmt.Sprintf("%ss", key)
	if value, ok := items[keys]; ok {
		if values, err := SliceOf[T](value); err == nil {
			result = append(result, values...)
		}
	}
	return result
}
