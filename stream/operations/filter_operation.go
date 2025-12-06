package operations

import (
	"gostream/stream/delegates"
	"gostream/stream/iterators"
)

type FilterOperation[T any] struct {
	predicate delegates.Predicate[T]
}

func AsFilterOperation[T any](filter delegates.Predicate[T]) *FilterOperation[T] {
	return &FilterOperation[T]{
		predicate: filter,
	}
}
func (f FilterOperation[T]) CanFuse(next StreamOperation[T]) bool {
	return true
}

func (f FilterOperation[T]) Apply(source iterators.Iterator[T]) iterators.Iterator[T] {
	return iterators.AsFilterIterator[T](f.predicate, source)
}

func (f FilterOperation[T]) ApplyFunc(func(args ...interface{}) interface{}) iterators.Iterator[T] {
	//TODO implement me
	panic("implement me")
}
