package operations

import "github.com/raimialiu/gostream/stream/iterators"

type TakeOperation[T any] struct {
	count int
}

func AsTakeOperation[T any](count int) *TakeOperation[T] {
	return &TakeOperation[T]{count: count}
}

func (t TakeOperation[T]) CanFuse(next StreamOperation[T]) bool { return true }

func (t TakeOperation[T]) Apply(source iterators.Iterator[T]) iterators.Iterator[T] {
	return iterators.AsTakeIterator(source, t.count)
}

func (t TakeOperation[T]) ApplyFunc(func(args ...interface{}) interface{}) iterators.Iterator[T] {
	panic("not implemented")
}
