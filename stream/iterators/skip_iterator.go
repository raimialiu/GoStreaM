package iterators

type SkipIterator[T any] struct {
	source  Iterator[T]
	count   int
	skipped bool
}

func AsSkipIterator[T any](source Iterator[T], count int) *SkipIterator[T] {
	return &SkipIterator[T]{
		source: source,
		count:  count,
	}
}

func (it *SkipIterator[T]) HasNext() bool {
	it.skip()
	return it.source.HasNext()
}

func (it *SkipIterator[T]) Next() T {
	it.skip()
	if !it.source.HasNext() {
		var zero T
		return zero
	}
	return it.source.Next()
}

func (it *SkipIterator[T]) skip() {
	if it.skipped {
		return
	}
	it.skipped = true
	for i := 0; i < it.count && it.source.HasNext(); i++ {
		it.source.Next()
	}
}

func (it *SkipIterator[T]) Close() error {
	return it.source.Close()
}
