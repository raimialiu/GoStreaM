package iterators

type TakeWhileIterator[T any] struct {
	source    Iterator[T]
	predicate func(T) bool
	next      T
	hasNext   bool
	computed  bool
	done      bool
}

func AsTakeWhileIterator[T any](source Iterator[T], predicate func(T) bool) *TakeWhileIterator[T] {
	return &TakeWhileIterator[T]{
		source:    source,
		predicate: predicate,
	}
}

func (it *TakeWhileIterator[T]) HasNext() bool {
	if !it.computed {
		it.computeNext()
	}
	return it.hasNext
}

func (it *TakeWhileIterator[T]) Next() T {
	if !it.HasNext() {
		var zero T
		return zero
	}
	result := it.next
	it.computed = false
	it.hasNext = false
	return result
}

func (it *TakeWhileIterator[T]) computeNext() {
	if it.done || !it.source.HasNext() {
		it.hasNext = false
		it.computed = true
		return
	}
	value := it.source.Next()
	if it.predicate(value) {
		it.next = value
		it.hasNext = true
	} else {
		it.done = true
		it.hasNext = false
	}
	it.computed = true
}

func (it *TakeWhileIterator[T]) Close() error {
	return it.source.Close()
}
