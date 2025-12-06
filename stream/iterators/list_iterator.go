package iterators

import "reflect"

type ListIterator[T any] struct {
	items        []T
	currentIndex int
	item         T
	zero         T
}

func (it *ListIterator[T]) Close() error {
	it.items = make([]T, 0)
	it.currentIndex = 0
	return nil
}

func AsListIterator[T any](items ...T) *ListIterator[T] {
	var zero T
	return &ListIterator[T]{
		items:        items,
		currentIndex: 0,
		zero:         zero,
	}
}

func (it *ListIterator[T]) HasNext() bool {
	return len(it.items) > it.currentIndex || !reflect.DeepEqual(it.item, it.zero)
}

func (it *ListIterator[T]) Next() T {
	if !it.HasNext() {
		it.currentIndex = 0
		var zero T
		return zero
	}

	value := it.items[it.currentIndex]
	it.currentIndex++
	return value
}
