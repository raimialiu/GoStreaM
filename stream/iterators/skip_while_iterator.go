package iterators

type SkipWhileIterator[T any] struct {
	source    Iterator[T]
	predicate func(T) bool
	next      T
	hasNext   bool
	computed  bool
	skipping  bool
}

func AsSkipWhileIterator[T any](source Iterator[T], predicate func(T) bool) *SkipWhileIterator[T] {
	return &SkipWhileIterator[T]{
		source:    source,
		predicate: predicate,
		skipping:  true,
	}
}

func (it *SkipWhileIterator[T]) HasNext() bool {
	if !it.computed {
		it.computeNext()
	}
	return it.hasNext
}

func (it *SkipWhileIterator[T]) Next() T {
	if !it.HasNext() {
		var zero T
		return zero
	}
	result := it.next
	it.computed = false
	it.hasNext = false
	return result
}

func (it *SkipWhileIterator[T]) computeNext() {
	for it.source.HasNext() {
		value := it.source.Next()
		if it.skipping {
			if it.predicate(value) {
				continue
			}
			it.skipping = false
		}
		it.next = value
		it.hasNext = true
		it.computed = true
		return
	}
	it.hasNext = false
	it.computed = true
}

func (it *SkipWhileIterator[T]) Close() error {
	return it.source.Close()
}
