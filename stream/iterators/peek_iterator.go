package iterators

type PeekIterator[T any] struct {
	source Iterator[T]
	action func(T)
}

func AsPeekIterator[T any](source Iterator[T], action func(T)) *PeekIterator[T] {
	return &PeekIterator[T]{
		source: source,
		action: action,
	}
}

func (it *PeekIterator[T]) HasNext() bool {
	return it.source.HasNext()
}

func (it *PeekIterator[T]) Next() T {
	if !it.HasNext() {
		var zero T
		return zero
	}
	value := it.source.Next()
	it.action(value)
	return value
}

func (it *PeekIterator[T]) Close() error {
	return it.source.Close()
}
