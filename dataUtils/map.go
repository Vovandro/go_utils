package dataUtils

// MapKeys returns a slice of keys from the provided map.
func MapKeys[K comparable, V any](data map[K]V) []K {
	keys := make([]K, len(data))

	i := 0
	for k := range data {
		keys[i] = k
		i++
	}

	return keys
}

// MapValues returns a slice of values from the provided map.
func MapValues[K comparable, V any](data map[K]V) []V {
	values := make([]V, len(data))

	i := 0
	for _, v := range data {
		values[i] = v
		i++
	}

	return values
}

// MapContains checks if the provided map contains the specified key.
func MapContains[K comparable, V any](data map[K]V, key K) bool {
	_, ok := data[key]
	return ok
}

// MapValueContains checks if the provided map contains the specified value.
func MapValueContains[K, V comparable](data map[K]V, value V) bool {
	for _, v := range data {
		if v == value {
			return true
		}
	}

	return false
}

// MapOrderedIterator returns a function that iterates over the provided map in a sorted order.
func MapOrderedIterator[K Constraints, V any](data map[K]V) func(yield func(K, V) bool) {
	keys := MapKeys(data)
	SliceOrderAsk(&keys)

	return func(yield func(K, V) bool) {
		for _, k := range keys {
			if !yield(k, data[k]) {
				return
			}
		}
	}
}

func MapToSlice[K comparable, V, T any](data map[K]V, predicate func(K, V) T) []T {
	values := make([]T, len(data))

	i := 0
	for k, v := range data {
		values[i] = predicate(k, v)
		i++
	}

	return values
}
