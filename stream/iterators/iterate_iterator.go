package iterators

type iterateIterator[T any] struct {
	current T
	hasNext func(T) bool
	next    func(T) T
	first   bool
}

func AsIterateIterator[T any](seed T, hasNext func(T) bool, next func(T) T) *iterateIterator[T] {
	return &iterateIterator[T]{
		current: seed,
		hasNext: hasNext,
		next:    next,
		first:   true,
	}
}

func (it *iterateIterator[T]) HasNext() bool {
	return it.hasNext(it.current)
}

func (it *iterateIterator[T]) Next() T {
	if !it.HasNext() {
		var zero T
		return zero
	}

	if it.first {
		it.first = false
		return it.current
	}

	it.current = it.next(it.current)
	return it.current
}

func (it *iterateIterator[T]) Close() error {
	return nil
}
