package iterators

type iterateIterator[T any] struct {
	current     T
	hasNextFn   func(T) bool
	nextFn      func(T) T
	started     bool
	nextVal     T
	hasComputed bool
	hasMore     bool
}

func AsIterateIterator[T any](seed T, hasNext func(T) bool, next func(T) T) *iterateIterator[T] {
	return &iterateIterator[T]{
		current:   seed,
		hasNextFn: hasNext,
		nextFn:    next,
		started:   false,
	}
}

func (it *iterateIterator[T]) HasNext() bool {
	if !it.started {
		return it.hasNextFn(it.current)
	}
	if !it.hasComputed {
		it.nextVal = it.nextFn(it.current)
		it.hasMore = it.hasNextFn(it.nextVal)
		it.hasComputed = true
	}
	return it.hasMore
}

func (it *iterateIterator[T]) Next() T {
	if !it.HasNext() {
		var zero T
		return zero
	}

	if !it.started {
		it.started = true
		return it.current
	}

	it.current = it.nextVal
	it.hasComputed = false
	return it.current
}

func (it *iterateIterator[T]) Close() error {
	return nil
}
