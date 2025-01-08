package dataUtils

func SliceDistinct[T comparable](list []T) []T {
	unique := make(map[T]bool)
	result := make([]T, 0, len(list))

	for _, item := range list {
		if !unique[item] {
			unique[item] = true
			result = append(result, item)
		}
	}

	return result
}

func SliceFilter[T any](slice []T, predicate func(*T) bool) []T {
	result := make([]T, 0, len(slice))

	for _, v := range slice {
		if predicate(&v) {
			result = append(result, v)
		}
	}
	return result
}

func SliceForeach[T any](slice *[]T, predicate func(*T)) {
	for i := range *slice {
		predicate(&(*slice)[i])
	}
}

func SliceToMap[T any, K comparable, V any](slice *[]T, predicate func(*T) (K, V)) map[K]V {
	result := make(map[K]V)
	for i := range *slice {
		k, v := predicate(&(*slice)[i])
		result[k] = v
	}

	return result
}
