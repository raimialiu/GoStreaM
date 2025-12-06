package iterators

type EmptyIterator[T any] struct {
	items        []T
	currentIndex int
}

func AsEmptyIterator[T any]() *EmptyIterator[T] {
	return &EmptyIterator[T]{
		items:        make([]T, 0),
		currentIndex: 0,
	}
}

func (it *EmptyIterator[T]) Next() T {
	if !it.HasNext() {
		var zero T
		return zero
	}

	value := it.items[it.currentIndex]
	it.currentIndex++
	it.items = it.items[it.currentIndex:]
	if len(it.items) == 0 {
		it.Close()
	}
	return value
}

func (it *EmptyIterator[T]) HasNext() bool {
	return it.currentIndex < len(it.items)
}

func (it *EmptyIterator[T]) Close() error {
	it.items = make([]T, 0)
	it.currentIndex = 0
	return nil
}
