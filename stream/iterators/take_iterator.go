package iterators

type TakeIterator[T any] struct {
	source    Iterator[T]
	remaining int
}

func AsTakeIterator[T any](source Iterator[T], count int) *TakeIterator[T] {
	return &TakeIterator[T]{
		source:    source,
		remaining: count,
	}
}

func (it *TakeIterator[T]) HasNext() bool {
	return it.remaining > 0 && it.source.HasNext()
}

func (it *TakeIterator[T]) Next() T {
	if !it.HasNext() {
		var zero T
		return zero
	}
	it.remaining--
	return it.source.Next()
}

func (it *TakeIterator[T]) Close() error {
	return it.source.Close()
}
