package iterators

type ConcatIterator[T any] struct {
	iterators []Iterator[T]
	index     int
}

func AsConcatIterator[T any](iterators ...Iterator[T]) *ConcatIterator[T] {
	return &ConcatIterator[T]{
		iterators: iterators,
		index:     0,
	}
}

func (it *ConcatIterator[T]) HasNext() bool {
	for it.index < len(it.iterators) {
		if it.iterators[it.index].HasNext() {
			return true
		}
		it.index++
	}
	return false
}

func (it *ConcatIterator[T]) Next() T {
	if !it.HasNext() {
		var zero T
		return zero
	}
	return it.iterators[it.index].Next()
}

func (it *ConcatIterator[T]) Close() error {
	for _, iter := range it.iterators {
		iter.Close()
	}
	return nil
}
