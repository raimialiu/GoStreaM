package iterators

type InfiniteIterator[T any] struct {
	value  T
	closed bool
}

func AsInfiniteIterator[T any](value T) *InfiniteIterator[T] {
	return &InfiniteIterator[T]{value: value, closed: false}
}
func (it *InfiniteIterator[T]) Next() T {
	if !it.HasNext() {
		var zero T
		return zero
	}
	return it.value
}

func (it *InfiniteIterator[T]) HasNext() bool {
	return it.closed
}

func (it *InfiniteIterator[T]) Close() error {
	var zero T
	it.value = zero
	it.closed = true
	return nil
}
