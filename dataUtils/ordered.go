package dataUtils

import "sort"

// Constraints is a type constraint that defines the types that can be used with the SliceOrderAsk and SliceOrderDesc functions.
type Constraints interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64 | ~string
}

// SliceOrderAsk sorts the provided slice in ascending order.
func SliceOrderAsk[T Constraints](slice *[]T) {
	sort.Slice(*slice, func(i, j int) bool {
		return (*slice)[i] < (*slice)[j]
	})
}

// SliceOrderDesc sorts the provided slice in descending order.
func SliceOrderDesc[T Constraints](slice *[]T) {
	sort.Slice(*slice, func(i, j int) bool {
		return (*slice)[i] > (*slice)[j]
	})
}
