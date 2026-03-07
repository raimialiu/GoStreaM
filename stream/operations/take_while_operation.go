package operations

import "github.com/raimialiu/gostream/stream/iterators"

type TakeWhileOperation[T any] struct {
	predicate func(T) bool
}

func AsTakeWhileOperation[T any](predicate func(T) bool) *TakeWhileOperation[T] {
	return &TakeWhileOperation[T]{predicate: predicate}
}

func (t TakeWhileOperation[T]) CanFuse(next StreamOperation[T]) bool { return true }

func (t TakeWhileOperation[T]) Apply(source iterators.Iterator[T]) iterators.Iterator[T] {
	return iterators.AsTakeWhileIterator(source, t.predicate)
}

func (t TakeWhileOperation[T]) ApplyFunc(func(args ...interface{}) interface{}) iterators.Iterator[T] {
	panic("not implemented")
}
