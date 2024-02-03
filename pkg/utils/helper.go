package utils

func MapFunc[T, U any](arr []T, fn func(T) U) []U {
	mappedArr := make([]U, len(arr))
	for i, v := range arr {
		mappedArr[i] = fn(v)
	}
	return mappedArr
}
