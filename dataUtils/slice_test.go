package dataUtils

import (
	"reflect"
	"testing"
)

func TestSliceDistinct(t *testing.T) {
	t.Run("test slice distinct int", func(t *testing.T) {
		got := SliceDistinct([]int{1, 1, 5, 2, 1, 5})
		SliceOrderAsk(&got)
		if !reflect.DeepEqual(got, []int{1, 2, 5}) {
			t.Errorf("SliceDistinct() = %v, want %v", got, []int{1, 2, 5})
		}
	})

	t.Run("test slice distinct string", func(t *testing.T) {
		got := SliceDistinct([]string{"1", "1", "5", "2", "1", "5"})
		SliceOrderDesc(&got)
		if !reflect.DeepEqual(got, []string{"5", "2", "1"}) {
			t.Errorf("SliceDistinct() = %v, want %v", got, []string{"5", "2", "1"})
		}
	})
}

func TestSliceFilter(t *testing.T) {
	t.Run("test slice filter int", func(t *testing.T) {
		if got := SliceFilter([]int{1, 1, 5, 2, 1, 5}, func(v *int) bool { return *v > 1 }); !reflect.DeepEqual(got, []int{5, 2, 5}) {
			t.Errorf("SliceFilter() = %v, want %v", got, []int{5, 2, 5})
		}
	})
}

func TestSliceForeach(t *testing.T) {
	t.Run("test slice foreach int", func(t *testing.T) {
		slice := []int{1, 2, 3}
		SliceForeach(&slice, func(v *int) {
			*v *= 2
		})
		if !reflect.DeepEqual(slice, []int{2, 4, 6}) {
			t.Errorf("SliceForeach() = %v, want %v", slice, []int{2, 4, 6})
		}
	})
}
