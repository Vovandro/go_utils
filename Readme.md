
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