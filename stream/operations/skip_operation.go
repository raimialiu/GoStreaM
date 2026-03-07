package operations

import "github.com/raimialiu/gostream/stream/iterators"

type SkipOperation[T any] struct {
	count int
}

func AsSkipOperation[T any](count int) *SkipOperation[T] {
	return &SkipOperation[T]{count: count}
}

func (s SkipOperation[T]) CanFuse(next StreamOperation[T]) bool { return true }

func (s SkipOperation[T]) Apply(source iterators.Iterator[T]) iterators.Iterator[T] {
	return iterators.AsSkipIterator(source, s.count)
}

func (s SkipOperation[T]) ApplyFunc(func(args ...interface{}) interface{}) iterators.Iterator[T] {
	panic("not implemented")
}
