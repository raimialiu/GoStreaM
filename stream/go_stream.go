package stream

import (
	"context"
	"gostream/stream/iterators"
	operations2 "gostream/stream/operations"
	"gostream/stream/types"
	"io"
	"sync"
)

type (
	GoStream[T any] struct {
		iterator   iterators.Iterator[T]
		operations []operations2.StreamOperation[T]

		// parallel indicates if this stream should use parallel processing
		parallel bool

		// closed indicates if the stream has been closed
		closed bool

		mu sync.Mutex

		ctx context.Context
	}
)

func (s *GoStream[T]) withOperations(operations ...operations2.StreamOperation[T]) *GoStream[T] {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.operations = append(s.operations, operations...)
	return s
}

func (s *GoStream[T]) withOperation(operations operations2.StreamOperation[T]) *GoStream[T] {
	return s.withOperations(operations)
}

func newGoStream[T any](iterator iterators.Iterator[T]) *GoStream[T] {
	return &GoStream[T]{
		iterator: iterator,
		parallel: false,
		closed:   false,
	}
}

func From[T any](slice []T) *GoStream[T] {
	return newGoStream(iterators.AsListIterator[T](slice...))
}

func Of[T any](items ...T) *GoStream[T] {
	return newGoStream(iterators.AsListIterator[T](items...))
}

func Empty[T any]() *GoStream[T] {
	return newGoStream(iterators.AsEmptyIterator[T]())
}

func Range[T types.Numeric](start, stop, step int) *GoStream[T] {
	return newGoStream(iterators.AsRangeIterator[T](start, stop, step))
}

func LongRange(start, stop, step int) *GoStream[int64] {
	return newGoStream(iterators.AsRangeIterator[int64](start, stop, step))
}

func IntRange(start, stop, step int) *GoStream[int] {
	return newGoStream(iterators.AsRangeIterator[int](start, stop, step))
}

func DoubleRange(start, stop, step float32) *GoStream[float32] {
	return newGoStream(iterators.AsDoubleRangeIterator(start, stop, step))
}

func Repeat[T any](value T, count int) *GoStream[T] {
	return newGoStream(iterators.AsRepeatIterator[T](value, count))
}

func RepeatForever[T any](value T) *GoStream[T] {
	return newGoStream(iterators.AsInfiniteIterator(value))
}

func FromReader[T any](r io.Reader, bufferSize int) *GoStream[string] {
	return newGoStream(iterators.AsReaderIterator[string](r, bufferSize))
}

func Generate[T any](generator func(args ...interface{}) T) *GoStream[T] {
	return newGoStream(iterators.AsGeneratorIterator(generator))
}

func Iterate[T any](seed T, hasNext func(T) bool, next func(T) T) *GoStream[T] {
	return newGoStream(iterators.AsIterateIterator(seed, hasNext, next))
}

func FromMap[K comparable, V any](m map[K]V) *GoStream[KeyValue[K, V]] {
	pairs := make([]KeyValue[K, V], 0, len(m))
	for k, v := range m {
		pairs = append(pairs, KeyValue[K, V]{Key: k, Value: v})
	}
	return newGoStream(iterators.AsListIterator(pairs...))
}

func FromMapValues[K comparable, V any](m map[K]V) *GoStream[V] {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return newGoStream(iterators.AsListIterator(values...))
}

func FromChannel[T any](channel <-chan T) *GoStream[T] {
	return newGoStream(iterators.AsChannelIterator(channel))
}

func FromFile(path string) *GoStream[string] {
	return newGoStream(iterators.AsFileIterator(path))
}
