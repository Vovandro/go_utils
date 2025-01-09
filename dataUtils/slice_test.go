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

func TestSliceToMap(t *testing.T) {
	// Тестовые данные
	tests := []struct {
		name      string
		input     []int
		predicate func(*int) (int, int)
		expected  map[int]int
	}{
		{
			name:  "Basic conversion of integers to their squares",
			input: []int{1, 2, 3, 4},
			predicate: func(i *int) (int, int) {
				return *i, (*i) * (*i)
			},
			expected: map[int]int{1: 1, 2: 4, 3: 9, 4: 16},
		},
		{
			name:  "Empty slice",
			input: []int{},
			predicate: func(i *int) (int, int) {
				return *i, (*i) * (*i)
			},
			expected: map[int]int{},
		},
		{
			name:  "Slice with duplicate keys",
			input: []int{1, 2, 2, 3},
			predicate: func(i *int) (int, int) {
				return *i, (*i) * (*i)
			},
			expected: map[int]int{1: 1, 2: 4, 3: 9}, // Последнее значение для дублирующегося ключа
		},
	}

	// Проход по тестам
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SliceToMap(&tt.input, tt.predicate)

			// Проверка результата
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SliceToMap() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
