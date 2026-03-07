package operations

import "github.com/raimialiu/gostream/stream/iterators"

type DistinctOperation[T any] struct {
	source iterators.Iterator[T]
}

func (d DistinctOperation[T]) CanFuse(next StreamOperation[T]) bool {
	return true
}

func (d DistinctOperation[T]) Apply(source iterators.Iterator[T]) iterators.Iterator[T] {
	return iterators.AsDistinctIterator[T](source)
}

func AsDistinct[T any]() *DistinctOperation[T] {
	return &DistinctOperation[T]{}
}
