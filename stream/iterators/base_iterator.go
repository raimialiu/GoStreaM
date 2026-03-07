package iterators

type BaseIterator[T any] struct {
	items        []T
	currentIndex int
}

func AsBaseIterator[T any](items ...T) *BaseIterator[T] {
	return &BaseIterator[T]{
		items:        items,
		currentIndex: 0,
	}
}

func (it *BaseIterator[T]) Add(items ...T) *BaseIterator[T] {
	it.items = append(it.items, items...)
	return it
}

func (it *BaseIterator[T]) Zero() T {
	var zero T
	return zero
}

func (it *BaseIterator[T]) HasNext() bool {
	return it.currentIndex < len(it.items)
}

func (it *BaseIterator[T]) Next() T {
	if !it.HasNext() {
		var zero T
		return zero
	}

	value := it.items[it.currentIndex]
	it.currentIndex++
	return value
}

func (it *BaseIterator[T]) Close() error {
	it.items = nil
	it.currentIndex = 0
	return nil
}
