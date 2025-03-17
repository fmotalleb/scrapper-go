package utils

func FirstOr[T any](arr []T, def T) T {
	if len(arr) != 0 {
		return arr[0]
	}
	return def
}

func FirstOrLazy[T any](arr []T, def func() T) T {
	if len(arr) != 0 {
		return arr[0]
	}
	return def()
}
