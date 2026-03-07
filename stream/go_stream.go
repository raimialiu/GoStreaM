package stream

import (
	"context"
	"github.com/raimialiu/gostream/stream/iterators"
	operations2 "github.com/raimialiu/gostream/stream/operations"
	"github.com/raimialiu/gostream/stream/types"
	"io"
	"sync"
)

type GoStream[T any] struct {
	iterator   iterators.Iterator[T]
	iterators  []iterators.Iterator[T]
	operations []operations2.StreamOperation[T]
	parallel   bool
	closed     bool
	mu         sync.Mutex
	ctx        context.Context
}

func (g *GoStream[T]) WithContext(ctx context.Context) *GoStream[T] {
	g.ctx = ctx
	return g
}

func (g *GoStream[T]) Iterators() []iterators.Iterator[T] {
	return g.iterators
}

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
	iters := make([]iterators.Iterator[T], 0)
	iters = append(iters, iterator)
	return &GoStream[T]{
		iterator:  iterator,
		iterators: iters,
		parallel:  false,
		closed:    false,
	}
}

func From[T any](slice []T) *GoStream[T] {
	return newGoStream(iterators.AsListIterator[T](slice...))
}

func New[T any]() *GoStream[T] {
	return Empty[T]()
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

func Concat[T any](streams ...*GoStream[T]) *GoStream[T] {
	if len(streams) == 0 {
		return Empty[T]()
	}
	result := streams[0]
	for _, s := range streams[1:] {
		result.iterators = append(result.iterators, s.iterators...)
	}
	return result
}

func Chunk[T any](s *GoStream[T], size int) *GoStream[[]T] {
	items := s.collect()
	var chunks [][]T
	for i := 0; i < len(items); i += size {
		end := i + size
		if end > len(items) {
			end = len(items)
		}
		chunks = append(chunks, items[i:end])
	}
	return From(chunks)
}

func Window[T any](s *GoStream[T], size int) *GoStream[[]T] {
	items := s.collect()
	var windows [][]T
	for i := 0; i <= len(items)-size; i++ {
		window := make([]T, size)
		copy(window, items[i:i+size])
		windows = append(windows, window)
	}
	return From(windows)
}
