package operations

import "github.com/raimialiu/gostream/stream/iterators"

type SkipWhileOperation[T any] struct {
	predicate func(T) bool
}

func AsSkipWhileOperation[T any](predicate func(T) bool) *SkipWhileOperation[T] {
	return &SkipWhileOperation[T]{predicate: predicate}
}

func (s SkipWhileOperation[T]) CanFuse(next StreamOperation[T]) bool { return true }

func (s SkipWhileOperation[T]) Apply(source iterators.Iterator[T]) iterators.Iterator[T] {
	return iterators.AsSkipWhileIterator(source, s.predicate)
}
