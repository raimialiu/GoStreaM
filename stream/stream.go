package stream

import (
	"gostream/stream/delegates"
	"gostream/stream/iterators"
)

const (
	Version      = "1.0.0"
	MinGoVersion = "1.21"
)

type (
	KeyValue[K comparable, V any] struct {
		Key   K
		Value V
	}

	// Stream  represents a sequence of elements supporting sequential and parallel
	//
	/// aggregate operations. This is the main interface for the Stream library.
	Stream[T any] interface {

		// ========== INTERMEDIATE OPERATIONS (Lazy) ==========

		// Filter returns a stream with elements matching the predicate
		Filter(predicate delegates.Predicate[T]) Stream[T]

		// Map transforms each element using the mapper function
		Map(mapper func(T) interface{}) Stream[interface{}]

		// FlatMap transforms each element to a stream and flattens the result
		FlatMap(mapper func(T) Stream[interface{}]) Stream[interface{}]

		// Distinct returns a stream with unique elements
		Distinct() Stream[T]

		// DistinctBy returns a stream with unique elements based on key selector
		DistinctBy(keySelector func(T)) Stream[T]

		// Take returns a stream with at most count elements
		Take(count int) Stream[T]

		// Skip returns a stream with the first count elements removed
		Skip(count int) Stream[T]

		// TakeWhile returns elements while predicate is true
		TakeWhile(predicate func(T) bool) Stream[T]

		// SkipWhile skips elements while predicate is true
		SkipWhile(predicate func(T) bool) Stream[T]

		// OrderByDesc sorts elements by key in descending order
		OrderByDesc(keySelector func(T) interface{}) OrderedStream[T]

		// Reverse returns a stream with elements in reverse order
		Reverse() Stream[T]

		// Concat concatenates this stream with another
		Concat(other Stream[T]) Stream[T]

		// Union returns the set union with another stream
		Union(other Stream[T]) Stream[T]

		// Intersect returns the set intersection with another stream
		Intersect(other Stream[T]) Stream[T]

		// Except returns elements not in the other stream
		Except(other Stream[T]) Stream[T]

		// Peek performs an action on each element without modifying the stream
		Peek(action func(T)) Stream[T]

		// Cache materializes the stream for reuse
		Cache() Stream[T]

		// Parallel enables parallel processing
		Parallel() ParallelStream[T]

		// ========== TERMINAL OPERATIONS (Eager) ==========

		// ForEach performs an action on each element
		ForEach(action func(T))

		// ToSlice collects elements into a slice
		ToSlice() []T

		// ToMap collects elements into a map using key selector
		ToMap(keySelector func(T) interface{}) map[interface{}]interface{}

		// Count returns the number of elements
		Count() int64

		// Any returns true if the stream has any elements
		Any() bool

		// AnyMatch returns true if any element matches the predicate
		AnyMatch(predicate func(T) bool) bool

		// All returns true if all elements match the predicate
		All(predicate func(T) bool) bool

		// None returns true if no elements match the predicate
		None(predicate func(T) bool) bool

		// First returns the first element or error if empty
		First() (T, error)

		// FirstOrDefault returns the first element or default value
		FirstOrDefault(defaultValue T) T

		// Last returns the last element or error if empty
		Last() (T, error)

		// Single returns the single element or error
		Single() (T, error)

		// Reduce reduces elements using accumulator function
		Reduce(identity T, accumulator func(T, T) T) T

		// ========== UTILITY METHODS ==========

		// Iterator returns an iterator for this stream
		Iterator() iterators.Iterator[T]

		// IsParallel returns true if this is a parallel stream
		IsParallel() bool

		// Close releases any resources held by the stream
		Close() error
	}

	OrderedStream[T any] interface {
		Stream[T]

		// OrderBy sorts elements by key in ascending order
		OrderBy(selector delegates.KeySelector[T]) OrderedStream[T]
		OrderByFunc(selector delegates.KeySelector[T]) OrderedStream[T]
		ThenBy(selector delegates.KeySelector[T]) OrderedStream[T]
		ThenByDesc(selector delegates.KeySelector[T]) OrderedStream[T]
	}

	ParallelStream[T any] interface {
		Stream[T]
		WithWorker(count int) ParallelStream[T]
		Sequential() Stream[T]
	}

	IntStream interface {
		Stream[int]
	}

	DoubleStream interface {
		Stream[float32]
	}
)
