package iterators

type RepeatIterator[T any] struct {
	value T
	count int
}

func AsRepeatIterator[T any](item T, count int) *RepeatIterator[T] {
	return &RepeatIterator[T]{
		value: item,
		count: count,
	}
}

func (it *RepeatIterator[T]) HasNext() bool {
	return it.count > 0
}

func (it *RepeatIterator[T]) Next() T {
	if !it.HasNext() {
		var zero T
		return zero
	}

	it.count--
	return it.value
}

func (it *RepeatIterator[T]) Close() error {
	it.count = 0
	return nil
}
