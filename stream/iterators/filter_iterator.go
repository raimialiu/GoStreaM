package iterators

import "gostream/stream/delegates"

type FilterIterator[T any] struct {
	predicate delegates.Predicate[T]
	source    Iterator[T]
	next      T
	done      bool
	hasNext   bool
	computed  bool
}

func AsFilterIterator[T any](predicate delegates.Predicate[T], source Iterator[T]) *FilterIterator[T] {
	return &FilterIterator[T]{
		predicate: predicate,
		source:    source,
		computed:  false,
		hasNext:   false,
	}
}

func (it *FilterIterator[T]) HasNext() bool {
	if !it.hasNext {
		it.computeNext()
	}
	return it.hasNext
}

func (it *FilterIterator[T]) Next() T {
	if !it.HasNext() {
		var zero T
		return zero
	}

	result := it.next
	it.computed = false
	return result
}

func (it *FilterIterator[T]) computeNext() {
	for it.source.HasNext() {
		candidate := it.source.Next()
		if it.predicate(candidate) {
			it.next = candidate
			it.hasNext = true
			it.computed = true
			return
		}
	}
	it.hasNext = false
	it.computed = true
}

func (it *FilterIterator[T]) Close() error {
	return it.source.Close()
}
