package operations

import "github.com/raimialiu/gostream/stream/iterators"

type DistinctByOperation[T any] struct {
	fieldFunc func(T) interface{}
}

func AsDistinctBy[T any](fieldFunc func(T) interface{}) *DistinctByOperation[T] {
	return &DistinctByOperation[T]{
		fieldFunc: fieldFunc,
	}
}

func (d DistinctByOperation[T]) CanFuse(next StreamOperation[T]) bool {
	return true
}

func (d DistinctByOperation[T]) Apply(source iterators.Iterator[T]) iterators.Iterator[T] {
	return iterators.AsDistinctByIterator(source, d.fieldFunc)
}
