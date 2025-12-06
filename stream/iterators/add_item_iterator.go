package iterators

type AddItemIterator[T any] struct {
	iterator *Iterator[T]
	base     *BaseIterator[T]
	items    []T
}

func (a AddItemIterator[T]) HasNext() bool {
	return a.base.HasNext()
}

func (a AddItemIterator[T]) Next() T {
	if a.HasNext() {
		return a.base.Next()
	}

	return a.base.Zero()
}

func (a AddItemIterator[T]) Close() error {
	return a.base.Close()
}

func AsAddItemIterator[T any](items ...T) *AddItemIterator[T] {
	return &AddItemIterator[T]{
		base:  AsBaseIterator[T](items...),
		items: items,
	}
}
