package utils

func FirstOr[T any](arr []T, def T) T {
	if len(arr) != 0 {
		return arr[0]
	}
	return def
}
