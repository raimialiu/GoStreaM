package operations

import "github.com/raimialiu/gostream/stream/iterators"

type PeekOperation[T any] struct {
	action func(T)
}

func AsPeekOperation[T any](action func(T)) *PeekOperation[T] {
	return &PeekOperation[T]{action: action}
}

func (p PeekOperation[T]) CanFuse(next StreamOperation[T]) bool { return true }

func (p PeekOperation[T]) Apply(source iterators.Iterator[T]) iterators.Iterator[T] {
	return iterators.AsPeekIterator(source, p.action)
}

func (p PeekOperation[T]) ApplyFunc(func(args ...interface{}) interface{}) iterators.Iterator[T] {
	panic("not implemented")
}
