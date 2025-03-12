package utils

import "fmt"

func SliceOf[T any](item any) ([]T, error) {
	if value, ok := item.(T); ok {
		return []T{value}, nil
	} else if value, ok := item.([]T); ok {
		return value, nil
	}
	return nil, fmt.Errorf("failed to read slice from %v", item)
}
