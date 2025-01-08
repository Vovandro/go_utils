package dataUtils

import (
	"reflect"
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
