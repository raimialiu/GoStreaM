package iterators

type ListIterator[T any] struct {
	items        []T
	currentIndex int
}

func AsListIterator[T any](items ...T) *ListIterator[T] {
	return &ListIterator[T]{
		items:        items,
		currentIndex: 0,
	}
}

func (it *ListIterator[T]) HasNext() bool {
	return it.currentIndex < len(it.items)
}

func (it *ListIterator[T]) Next() T {
	if !it.HasNext() {
		var zero T
		return zero
	}

	value := it.items[it.currentIndex]
	it.currentIndex++
	return value
}

func (it *ListIterator[T]) Close() error {
	it.items = nil
	it.currentIndex = 0
	return nil
}
