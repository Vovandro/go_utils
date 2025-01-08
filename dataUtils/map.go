package dataUtils

func MapKeys[K comparable, V any](data map[K]V) []K {
	keys := make([]K, len(data))

	i := 0
	for k := range data {
		keys[i] = k
		i++
	}

	return keys
}

func MapValues[K comparable, V any](data map[K]V) []V {
	values := make([]V, len(data))

	i := 0
	for _, v := range data {
		values[i] = v
		i++
	}

	return values
}

func MapContains[K comparable, V any](data map[K]V, key K) bool {
	_, ok := data[key]
	return ok
}

func MapValueContains[K, V comparable](data map[K]V, value V) bool {
	for _, v := range data {
		if v == value {
			return true
		}
	}

	return false
}

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
