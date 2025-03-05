package concurrency

type result[T any] struct {
	err error
	val T
}

type Future[T any] struct {
	ch chan result[T]
}

func NewFuture[T any](action func() (T, error)) *Future[T] {
	f := &Future[T]{make(chan result[T])}
	go func() {
		val, err := action()
		f.ch <- result[T]{val: val, err: err}
		close(f.ch)
	}()

	return f
}

func (f *Future[T]) Get() (T, error) {
	res := <-f.ch
	return res.val, res.err
}
