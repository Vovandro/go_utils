package dataUtils

import (
	"reflect"
	"sort"
	"strconv"
	"testing"
)

func TestMapContains(t *testing.T) {
	t.Run("test contains string", func(t *testing.T) {
		if got := MapContains(map[string]int{"a": 1}, "a"); got != true {
			t.Errorf("MapContains() = %v, want %v", got, true)
		}
	})

	t.Run("test contains int", func(t *testing.T) {
		if got := MapContains(map[int]int{5: 1}, 5); got != true {
			t.Errorf("MapContains() = %v, want %v", got, true)
		}
	})

	t.Run("test contains fail", func(t *testing.T) {
		if got := MapContains(map[string]int{"a": 1}, "b"); got != false {
			t.Errorf("MapContains() = %v, want %v", got, false)
		}
	})
}

func TestMapKeys(t *testing.T) {
	t.Run("test keys string", func(t *testing.T) {
		got := MapKeys(map[string]string{"a": "1", "b": "2"})
		SliceOrderAsk(&got)
		if !reflect.DeepEqual(got, []string{"a", "b"}) {
			t.Errorf("MapKeys() = %v, want %v", got, []string{"a", "b"})
		}
	})

	t.Run("test keys float", func(t *testing.T) {
		got := MapKeys(map[float32]string{1.1: "1", 1.2: "2"})
		SliceOrderDesc(&got)
		if !reflect.DeepEqual(got, []float32{1.2, 1.1}) {
			t.Errorf("MapKeys() = %v, want %v", got, []float32{1.1, 1.2})
		}
	})
}

func TestMapValueContains(t *testing.T) {
	t.Run("test value contains", func(t *testing.T) {
		if got := MapValueContains(map[string]int{"a": 1}, 1); got != true {
			t.Errorf("MapValueContains() = %v, want %v", got, true)
		}
	})

	t.Run("test value not contains", func(t *testing.T) {
		if got := MapValueContains(map[string]int{"a": 1}, 2); got != false {
			t.Errorf("MapValueContains() = %v, want %v", got, false)
		}
	})
}

func TestMapValues(t *testing.T) {
	t.Run("test map values", func(t *testing.T) {
		got := MapValues(map[string]int{"a": 1, "b": 2})
		SliceOrderAsk(&got)
		if !reflect.DeepEqual(got, []int{1, 2}) {
			t.Errorf("MapValues() = %v, want %v", got, []int{1, 2})
		}
	})
}

func TestMapOrderedIterator(t *testing.T) {
	t.Run("test map ordered iterator", func(t *testing.T) {
		got := make([]string, 0)

		for k, _ := range MapOrderedIterator(map[string]int{"a": 1, "c": 2, "b": 3, "f": 5, "g": 6}) {
			if k == "f" {
				break
			}
			got = append(got, k)
		}

		if !reflect.DeepEqual(got, []string{"a", "b", "c"}) {
			t.Errorf("MapOrderedIterator() = %v, want %v", got, []string{"a", "b", "c"})
		}
	})
}

func TestMapToSlice(t *testing.T) {
	t.Run("test map to slice with string keys", func(t *testing.T) {
		data := map[string]int{"a": 1, "b": 2, "c": 3}
		got := MapToSlice(data, func(k string, v int) string {
			return k + ":" + strconv.Itoa(v)
		})

		// Sort for consistent comparison
		sort.Strings(got)
		expected := []string{"a:1", "b:2", "c:3"}
		sort.Strings(expected)

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("MapToSlice() = %v, want %v", got, expected)
		}
	})

	t.Run("test map to slice with int keys", func(t *testing.T) {
		data := map[int]string{1: "one", 2: "two", 3: "three"}
		got := MapToSlice(data, func(k int, v string) int {
			return k * len(v)
		})

		// Sort for consistent comparison
		sort.Ints(got)
		expected := []int{3, 6, 15} // 1*3, 2*3, 3*5
		sort.Ints(expected)

		if !reflect.DeepEqual(got, expected) {
			t.Errorf("MapToSlice() = %v, want %v", got, expected)
		}
	})

	t.Run("test map to slice with empty map", func(t *testing.T) {
		data := map[string]int{}
		got := MapToSlice(data, func(k string, v int) string {
			return k + ":" + strconv.Itoa(v)
		})

		if len(got) != 0 {
			t.Errorf("MapToSlice() = %v, want empty slice", got)
		}
	})
}
