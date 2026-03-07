package operations

import "github.com/raimialiu/gostream/stream/iterators"

type FlatMapOperation[T any] struct {
	mapper func(T) []T
}

func AsFlatMapOperation[T any](mapper func(T) []T) *FlatMapOperation[T] {
	return &FlatMapOperation[T]{mapper: mapper}
}

func (f FlatMapOperation[T]) CanFuse(next StreamOperation[T]) bool { return true }

func (f FlatMapOperation[T]) Apply(source iterators.Iterator[T]) iterators.Iterator[T] {
	return iterators.AsFlatMapIterator(source, f.mapper)
}
