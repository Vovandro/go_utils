
[![pipeline status](https://gitlab.com/devpro_studio/go_utils/badges/master/pipeline.svg)](https://gitlab.com/devpro_studio/go_utils/-/commits/master)
[![coverage report](https://gitlab.com/devpro_studio/go_utils/badges/master/coverage.svg)](https://gitlab.com/devpro_studio/go_utils/-/commits/master)
[![Latest Release](https://gitlab.com/devpro_studio/go_utils/-/badges/release.svg)](https://gitlab.com/devpro_studio/go_utils/-/releases)

# Golang helper utils


## Installation

Add the package to your project:

```bash
go get gitlab.com/devpro_studio/go_utils
```

Import it in your code:

```go
import "gitlab.com/devpro_studio/go_utils"
```

---

<details>
  <summary>Concurrency Utilities for Go</summary>

## Concurrency Utilities for Go

This package provides a set of utilities for working with channels and concurrency in Go. It simplifies the creation of pipelines, fan-in/fan-out patterns, batching, and parallel processing. Each function is generic and works with any data type.

---

### Features

- **Pipeline**: Apply a transformation function to a stream of input values.
- **Split**: Split a stream of input values into two output streams based on a function.
- **FanIn**: Merge multiple input streams into a single output stream.
- **FanOut**: Distribute values from an input stream to multiple output streams.
- **Batch**: Group input values into fixed-size batches.
- **Parallel**: Process input values concurrently using a worker pool.

---

### Usage

#### 1. Pipeline

Transforms a stream of input values using a provided function.

```go
input := make(chan int)
output := concurrency.Pipeline(func(i int) int { return i * 2 })(input)

go func() {
    for i := 1; i <= 5; i++ {
        input <- i
    }
    close(input)
}()

for val := range output {
    fmt.Println(val) // Output: 2, 4, 6, 8, 10
}
```

---

#### 2. Split

Splits a stream of input values into two streams based on a transformation function.

```go
input := make(chan int)
out1, out2 := concurrency.Split(input, func(i int) (int, string) {
    return i, fmt.Sprintf("Value: %d", i)
})

go func() {
    for i := 1; i <= 3; i++ {
        input <- i
    }
    close(input)
}()

for val := range out1 {
    fmt.Println(val) // Output: 1, 2, 3
}

for val := range out2 {
    fmt.Println(val) // Output: "Value: 1", "Value: 2", "Value: 3"
}
```

---

#### 3. FanIn

Merges multiple input streams into a single output stream.

```go
ch1 := make(chan int)
ch2 := make(chan int)
output := concurrency.FanIn(ch1, ch2)

go func() {
    ch1 <- 1
    ch1 <- 2
    close(ch1)
}()

go func() {
    ch2 <- 3
    ch2 <- 4
    close(ch2)
}()

for val := range output {
    fmt.Println(val) // Output: 1, 2, 3, 4 (order may vary)
}
```

---

#### 4. FanOut

Distributes values from a single input stream to multiple output streams.

```go
input := make(chan int)
outputs := concurrency.FanOut(input, 3, false)

go func() {
    for i := 1; i <= 5; i++ {
        input <- i
    }
    close(input)
}()

for i, out := range outputs {
    go func(index int, ch <-chan int) {
        for val := range ch {
            fmt.Printf("Output[%d]: %d\n", index, val)
        }
    }(i, out)
}
```

---

#### 5. Batch

Groups values from a stream into fixed-size batches.

```go
input := make(chan int)
output := concurrency.Batch(input, 2)

go func() {
    for i := 1; i <= 5; i++ {
        input <- i
    }
    close(input)
}()

for batch := range output {
    fmt.Println(batch) // Output: [1 2], [3 4], [5]
}
```

---

#### 6. Parallel

Processes values from a stream concurrently using a worker pool.

```go
input := make(chan int)
output := concurrency.Parallel(input, func(i int) int { return i * i }, 3)

go func() {
    for i := 1; i <= 5; i++ {
        input <- i
    }
    close(input)
}()

for val := range output {
    fmt.Println(val) // Output: 1, 4, 9, 16, 25 (order may vary)
}
```

</details>

---

<details>
  <summary>dataUtils Package</summary>

## dataUtils Package

The dataUtils package provides a collection of utility functions for working with Go maps and slices. It simplifies common operations such as retrieving keys or values, checking for the presence of keys or values, filtering, deduplication, sorting, and converting slices to maps. The package leverages Go generics and type constraints to ensure flexibility and type safety.

---

### Features

- **Retrieve Keys**: Extract all keys from a map.
- **Retrieve Values**: Extract all values from a map.
- **Check for Keys**: Verify if a key exists in a map.
- **Check for Values**: Verify if a value exists in a map.
- **Ordered Iteration**: Iterate over a map in sorted key order.
- **Generic Sorting**: Sort slices in ascending or descending order.
- **Deduplication**: Remove duplicate elements from a slice.
- **Filtering**: Extract elements that match a given condition.
- **Iteration**: Apply a function to each element in a slice.
- **Conversion**: Transform a slice into a map.

---

#### 1. MapKeys

Extracts all keys from the provided map.

```go
data := map[string]int{"a": 1, "b": 2, "c": 3}
keys := dataUtils.MapKeys(data)
fmt.Println(keys) // Output: [a b c] (order may vary)
```

---

#### 2. MapValues

Extracts all values from the provided map.

```go
data := map[string]int{"a": 1, "b": 2, "c": 3}
values := dataUtils.MapValues(data)
fmt.Println(values) // Output: [1 2 3] (order may vary)
```

---

#### 3. MapContains

Checks if the provided map contains the specified key.

```go
data := map[string]int{"a": 1, "b": 2}
exists := dataUtils.MapContains(data, "b")
fmt.Println(exists) // Output: true
```

---

#### 4. MapValueContains

Checks if the provided map contains the specified value.

```go
data := map[string]int{"a": 1, "b": 2}
exists := dataUtils.MapValueContains(data, 2)
fmt.Println(exists) // Output: true
```

---

#### 5. MapOrderedIterator

Creates a function to iterate over the provided map in sorted key order. Requires a Constraints type for sorting keys.

```go
data := map[int]string{3: "three", 1: "one", 2: "two"}

iterator := dataUtils.MapOrderedIterator(data)
iterator(func(k int, v string) bool {
    fmt.Printf("Key: %d, Value: %s\n", k, v)
    return true // Return false to stop iteration early
})
// Output:
// Key: 1, Value: one
// Key: 2, Value: two
// Key: 3, Value: three
```

---

#### 6. SliceOrder

Sorts a slice in ascending or descending order.

```go
numbers := []int{5, 3, 8, 1}
dataUtils.SliceOrderAsk(&numbers)
fmt.Println(numbers) // Output: [1 3 5 8]

strings := []string{"banana", "apple", "cherry"}
dataUtils.SliceOrderAsk(&strings)
fmt.Println(strings) // Output: [apple banana cherry]
```

```go
numbers := []int{5, 3, 8, 1}
dataUtils.SliceOrderDesc(&numbers)
fmt.Println(numbers) // Output: [8 5 3 1]

strings := []string{"banana", "apple", "cherry"}
dataUtils.SliceOrderDesc(&strings)
fmt.Println(strings) // Output: [cherry banana apple]
```

---

#### 7. SliceDistinct

Removes duplicate elements from a slice and returns a new slice with unique values.

```go
numbers := []int{1, 2, 2, 3, 4, 4}
uniqueNumbers := dataUtils.SliceDistinct(numbers)
fmt.Println(uniqueNumbers) // Output: [1 2 3 4]
```

---

#### 8. SliceFilter

Filters elements in a slice based on a predicate function.

```go
numbers := []int{1, 2, 3, 4, 5}
evenNumbers := dataUtils.SliceFilter(numbers, func(n *int) bool {
    return *n%2 == 0
})
fmt.Println(evenNumbers) // Output: [2 4]
```

---

#### 9. SliceForeach

Applies a function to each element in a slice.

```go
numbers := []int{1, 2, 3}
dataUtils.SliceForeach(&numbers, func(n *int) {
    *n *= 2
})
fmt.Println(numbers) // Output: [2 4 6]
```

---

#### 10. SliceToMap

Converts a slice to a map using a transformation function.

```go
people := []string{"Alice", "Bob", "Charlie"}
nameLengths := dataUtils.SliceToMap(&people, func(name *string) (string, int) {
    return *name, len(*name)
})
fmt.Println(nameLengths) // Output: map[Alice:5 Bob:3 Charlie:7]
```

</details>

---

<details>
  <summary>decode Package Overview</summary>

## decode Package Overview

The decode package provides a flexible utility for copying and transforming data between different Go data structures, such as maps and structs. It supports nested structures, type conversions, and field mapping using struct tags. This functionality is useful for scenarios like data serialization, deserialization, and mapping between different data representations.

---

### Key Features

#### Core Functionality
- **Data Transformation**: Copy and transform data between structs, maps, and slices.
- **Field Mapping**: Map fields between structs and maps using custom tags.
- **Nested Data Handling**: Recursively process nested data structures.

#### Error Handling
- **Type Mismatch**: Returns an error if source and destination types are incompatible.
- **Destination Validation**: Ensures the destination is a writable pointer.
- **Field Presence**: Optionally enforce strict checks for field presence in the destination.

#### Configuration Flags
The behavior of the `Decode` function can be customized using the following flags:
- **`DecoderStrongFoundDst`**: Enforces strict checks for destination field presence.
- **`DecoderStrongType`**: Ensures type safety and allows struct-to-map conversion.
- **`DecoderUnwrapStructToMap`**: Unwraps nested structs into maps for flexible data representation.

---

### Usage Example

#### Basic Struct-to-Struct Mapping
```go
type Source struct {
    Name  string `json:"name"`
    Age   int    `json:"age"`
}

type Destination struct {
    Name  string `json:"name"`
    Age   int    `json:"age"`
}

source := Source{Name: "Alice", Age: 30}
var destination Destination

err := decode.Decode(source, &destination, "json", decode.DecoderStrongFoundDst)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Println(destination) // Output: {Name: "Alice", Age: 30}
```

---

#### Struct-to-Map Conversion

```go
type Source struct {
    Name  string `json:"name"`
    Age   int    `json:"age"`
}

source := Source{Name: "Alice", Age: 30}
destination := make(map[string]interface{})

err := decode.Decode(source, &destination, "json", decode.DecoderUnwrapStructToMap)
if err != nil {
    log.Fatalf("Error: %v", err)
}

fmt.Println(destination) // Output: map[name: "Alice" age: 30]
```

</details>