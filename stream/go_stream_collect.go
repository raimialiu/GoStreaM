package stream

import (
	"github.com/raimialiu/gostream/stream/collectors"
	"github.com/raimialiu/gostream/stream/iterators"
	"github.com/raimialiu/gostream/stream/operations"
)

// Collect applies a Collector to the stream elements.
func Collect[T any, A any, R any](s *GoStream[T], collector collectors.Collector[T, A, R]) R {
	acc := collector.Supplier()
	for _, source := range s.execute() {
		for source.HasNext() {
			acc = collector.Accumulator(acc, source.Next())
		}
	}
	return collector.Finisher(acc)
}

// ========== Custom Extensibility ==========

// Apply adds a custom StreamOperation to the pipeline.
// This allows users to create their own lazy intermediate operations
// by implementing the operations.StreamOperation interface.
func (g *GoStream[T]) Apply(op operations.StreamOperation[T]) *GoStream[T] {
	g.withOperation(op)
	return g
}

// ApplyIterator wraps the stream's current pipeline with a custom iterator
// transformation. This is a lower-level extension point than Apply.
func (g *GoStream[T]) ApplyIterator(transform func(iterators.Iterator[T]) iterators.Iterator[T]) *GoStream[T] {
	g.withOperation(&funcOperation[T]{fn: transform})
	return g
}

// Transform applies a function that takes the stream and returns a new one.
// Useful for composing reusable pipeline fragments.
func (g *GoStream[T]) Transform(fn func(*GoStream[T]) *GoStream[T]) *GoStream[T] {
	return fn(g)
}

// FromIterator creates a GoStream from a custom Iterator implementation.
func FromIterator[T any](iter iterators.Iterator[T]) *GoStream[T] {
	return newGoStream(iter)
}

// ========== funcOperation (internal) ==========

type funcOperation[T any] struct {
	fn func(iterators.Iterator[T]) iterators.Iterator[T]
}

func (f *funcOperation[T]) CanFuse(next operations.StreamOperation[T]) bool { return true }

func (f *funcOperation[T]) Apply(source iterators.Iterator[T]) iterators.Iterator[T] {
	return f.fn(source)
}
