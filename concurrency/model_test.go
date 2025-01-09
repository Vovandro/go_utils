package concurrency

import (
	"reflect"
	"sort"
	"sync"
	"testing"
)

func TestPipeline(t *testing.T) {
	// Тестовые данные
	tests := []struct {
		name     string
		input    []int
		fn       func(int) int
		expected []int
	}{
		{
			name:  "Square numbers",
			input: []int{1, 2, 3, 4},
			fn: func(x int) int {
				return x * x
			},
			expected: []int{1, 4, 9, 16},
		},
		{
			name:  "Multiply by 2",
			input: []int{5, 10, 15},
			fn: func(x int) int {
				return x * 2
			},
			expected: []int{10, 20, 30},
		},
		{
			name:  "Empty input",
			input: []int{},
			fn: func(x int) int {
				return x + 1
			},
			expected: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем входной канал
			in := make(chan int, len(tt.input))
			for _, val := range tt.input {
				in <- val
			}
			close(in)

			// Применяем Pipeline
			pipeline := Pipeline(tt.fn)
			outChan := pipeline(in)

			// Считываем данные из выходного канала
			result := make([]int, 0)
			for val := range outChan {
				result = append(result, val)
			}

			// Проверяем результат
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Pipeline() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestFanIn(t *testing.T) {
	tests := []struct {
		name     string
		inputs   [][]int
		expected []int
	}{
		{"Merge two streams", [][]int{{1, 2, 3}, {4, 5}}, []int{1, 2, 3, 4, 5}},
		{"Empty streams", [][]int{{}, {}}, []int{}},
		{"Single stream", [][]int{{1, 2, 3}}, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем входные каналы
			var inputs []<-chan int
			for _, input := range tt.inputs {
				ch := make(chan int, len(input))
				for _, v := range input {
					ch <- v
				}
				close(ch)
				inputs = append(inputs, ch)
			}

			// Тестируем FanIn
			resultCh := FanIn(inputs...)
			result := make([]int, 0)
			for v := range resultCh {
				result = append(result, v)
			}

			sort.Ints(result)
			sort.Ints(tt.expected)

			// Проверяем результат
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestFanOut(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		size     int
		expected [][]int
	}{
		{"Distribute evenly", []int{1, 2, 3, 4, 5}, 2, [][]int{{1, 3, 5}, {2, 4}}},
		{"Single channel", []int{1, 2, 3}, 1, [][]int{{1, 2, 3}}},
		{"More channels than data", []int{1, 2}, 3, [][]int{{1}, {2}, {}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inCh := make(chan int, len(tt.input))
			for _, v := range tt.input {
				inCh <- v
			}
			close(inCh)

			var w sync.WaitGroup
			w.Add(tt.size)
			outChannels := FanOut(inCh, tt.size, true)
			resOut := make(chan struct {
				I    int
				Data []int
			}, tt.size)
			results := make([][]int, tt.size)

			for i, ch := range outChannels {
				go func(i int) {
					batch := make([]int, 0)
					for v := range ch {
						batch = append(batch, v)
					}
					resOut <- struct {
						I    int
						Data []int
					}{I: i, Data: batch}
					w.Done()
				}(i)
			}

			w.Wait()
			close(resOut)

			for v := range resOut {
				results[v.I] = v.Data
			}

			// Проверяем результат
			if len(results) != len(tt.expected) {
				t.Fatalf("expected %d channels, got %d", len(tt.expected), len(results))
			}

			for i := range results {
				sort.Ints(results[i])
				sort.Ints(tt.expected[i])

				if !reflect.DeepEqual(results[i], tt.expected[i]) {
					t.Errorf("channel %d: expected %v, got %v", i, tt.expected[i], results[i])
				}
			}
		})
	}
}

func TestBatch(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		size     int
		expected [][]int
	}{
		{"Exact batches", []int{1, 2, 3, 4, 5, 6}, 2, [][]int{{1, 2}, {3, 4}, {5, 6}}},
		{"Incomplete batch", []int{1, 2, 3, 4, 5}, 2, [][]int{{1, 2}, {3, 4}, {5}}},
		{"Single element batches", []int{1, 2, 3}, 1, [][]int{{1}, {2}, {3}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inCh := make(chan int, len(tt.input))
			for _, v := range tt.input {
				inCh <- v
			}
			close(inCh)

			resultCh := Batch(inCh, tt.size)

			var results [][]int
			for batch := range resultCh {
				results = append(results, batch)
			}

			// Проверяем результат
			if len(results) != len(tt.expected) {
				t.Fatalf("expected %d batches, got %d", len(tt.expected), len(results))
			}

			for i := range results {
				sort.Ints(results[i])
				sort.Ints(tt.expected[i])

				if !reflect.DeepEqual(results[i], tt.expected[i]) {
					t.Errorf("batch %d: expected %v, got %v", i, tt.expected[i], results[i])
				}
			}
		})
	}
}

func TestParallel(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		fn       func(int) int
		count    int
		expected []int
	}{
		{"Double values", []int{1, 2, 3, 4}, func(x int) int { return x * 2 }, 3, []int{2, 4, 6, 8}},
		{"Square values", []int{1, 2, 3}, func(x int) int { return x * x }, 2, []int{1, 4, 9}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inCh := make(chan int, len(tt.input))
			for _, v := range tt.input {
				inCh <- v
			}
			close(inCh)

			outCh := Parallel(inCh, tt.fn, tt.count)

			var results []int
			for v := range outCh {
				results = append(results, v)
			}

			sort.Ints(results)
			sort.Ints(tt.expected)

			if !reflect.DeepEqual(results, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, results)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		fn        func(int) (int, int)
		expected  []int
		expected2 []int
	}{
		{"split values", []int{12, 23, 34}, func(x int) (int, int) { return x / 10, x % 10 }, []int{1, 2, 3}, []int{2, 3, 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inCh := make(chan int, len(tt.input))
			for _, v := range tt.input {
				inCh <- v
			}
			close(inCh)

			got, got2 := Split(inCh, tt.fn)

			var results []int
			var results2 []int
			wg := &sync.WaitGroup{}
			wg.Add(2)

			go func() {
				for v := range got {
					results = append(results, v)
				}
				wg.Done()
			}()

			go func() {
				for v := range got2 {
					results2 = append(results2, v)
				}
				wg.Done()
			}()

			wg.Wait()

			sort.Ints(results)
			sort.Ints(results2)
			sort.Ints(tt.expected)
			sort.Ints(tt.expected2)

			if !reflect.DeepEqual(results, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, results)
			}

			if !reflect.DeepEqual(results2, tt.expected2) {
				t.Errorf("expected %v, got %v", tt.expected2, results2)
			}
		})
	}
}
