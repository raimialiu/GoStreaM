package iterators

type BaseIterator[T any] struct {
	items        []T
	currentIndex int
	item         T
	zero         T
}

func AsBaseIterator[T any](items ...T) *BaseIterator[T] {
	var zero T
	return &BaseIterator[T]{
		items:        items,
		currentIndex: 0,
		zero:         zero,
	}
}

func (it *BaseIterator[T]) Add(items ...T) *BaseIterator[T] {
	for _, v := range items {
		it.items = append(it.items, v)
	}

	return it
}

func (it *BaseIterator[T]) Zero() T {
	return it.zero
}

func (it *BaseIterator[T]) HasNext() bool {
	return len(it.items) > it.currentIndex
}

func (it *BaseIterator[T]) Next() T {
	if !it.HasNext() {
		it.currentIndex = 0
		var zero T
		it.item = zero
		return zero
	}

	value := it.items[it.currentIndex]
	it.item = value
	it.currentIndex++
	return value
}

func (it *BaseIterator[T]) Close() error {
	it.items = make([]T, 0)
	it.currentIndex = 0
	return nil
}
