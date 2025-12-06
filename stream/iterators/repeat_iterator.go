package iterators

type RepeatIterator[T any] struct {
	items []T // container to keep
	value T
	count int
}

func AsRepeatIterator[T any](item T, count int) *RepeatIterator[T] {
	return &RepeatIterator[T]{
		items: make([]T, 0),
		count: count,
	}
}

func (it *RepeatIterator[T]) HasNext() bool {
	canContinue := it.count > 0
	if !canContinue {
		it.Close()
	}

	return canContinue
}

func (it *RepeatIterator[T]) Next() T {
	if !it.HasNext() {
		if len(it.items) == 0 {
			var zero T
			return zero
		}

		return it.items[0]
	}

	it.count--
	value := it.items[it.count]

	return value
}

func (it *RepeatIterator[T]) Close() error {
	it.count = 0
	it.items = make([]T, 0)
	return nil
}

// 2, repeat (10)
