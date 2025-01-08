package dataUtils

import "sort"

type Constraints interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64 | ~string
}

func SliceOrderAsk[T Constraints](slice *[]T) {
	sort.Slice(*slice, func(i, j int) bool {
		return (*slice)[i] < (*slice)[j]
	})
}

func SliceOrderDesc[T Constraints](slice *[]T) {
	sort.Slice(*slice, func(i, j int) bool {
		return (*slice)[i] > (*slice)[j]
	})
}
