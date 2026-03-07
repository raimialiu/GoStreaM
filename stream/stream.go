package stream

import (
	"github.com/raimialiu/gostream/stream/delegates"
	"github.com/raimialiu/gostream/stream/iterators"
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

	Stream[T any] interface {
		// ========== INTERMEDIATE OPERATIONS (Lazy) ==========

		Filter(predicate delegates.Predicate[T]) *GoStream[T]
		Map(mapper func(T) interface{}) *GoStream[T]
		FlatMap(mapper func(T) []T) *GoStream[T]
		Distinct() *GoStream[T]
		DistinctBy(keySelector func(T) interface{}) *GoStream[T]
		Take(count int) *GoStream[T]
		Skip(count int) *GoStream[T]
		TakeWhile(predicate func(T) bool) *GoStream[T]
		SkipWhile(predicate func(T) bool) *GoStream[T]
		Peek(action func(T)) *GoStream[T]
		Reverse() *GoStream[T]
		Concat(others ...GoStream[T]) *GoStream[T]
		Union(other GoStream[T]) *GoStream[T]
		Intersect(other GoStream[T]) *GoStream[T]
		Except(other GoStream[T]) *GoStream[T]
		Cache() *GoStream[T]
		Select(selector func(T) interface{}) *GoStream[T]
		OrderBy(less func(a, b T) bool) *GoStream[T]
		OrderByDesc(less func(a, b T) bool) *GoStream[T]
		Zip(other *GoStream[T], zipper func(T, T) T) *GoStream[T]

		// ========== TERMINAL OPERATIONS (Eager) ==========

		ForEach(action func(T))
		ToList() []T
		ToSlice() []T
		ToMap(keySelector func(T) interface{}) map[interface{}]T
		GroupBy(keySelector func(T) interface{}) map[interface{}][]T
		Partition(predicate func(T) bool) ([]T, []T)
		Count() int
		CountBy(predicate func(T) bool) int
		Any() bool
		AnyMatch(predicate func(T) bool) bool
		AllMatch(predicate func(T) bool) bool
		NoneMatch(predicate func(T) bool) bool
		First() (T, error)
		FirstOrDefault(defaultValue T) T
		Last() (T, error)
		Single() (T, error)
		Reduce(identity T, accumulator func(T, T) T) T
		Contains(value T) bool
		Sum() float64
		Average() float64
		Min() (T, error)
		Max() (T, error)

		// ========== UTILITY METHODS ==========

		Iterator() iterators.Iterator[T]
		IsParallel() bool
		Close()
	}
)
