package operations

import "gostream/stream/iterators"

type StreamOperation[T any] interface {
	// CanFuse returns true if this operation can be fused with the next
	CanFuse(next StreamOperation[T]) bool

	Apply(source iterators.Iterator[T]) iterators.Iterator[T]
	ApplyFunc(func(args ...interface{}) interface{}) iterators.Iterator[T]
	// AsParallel() ParallelStream[T]
	/*
		// Collection terminals
		ToList() []T  // Alias for ToSlice
		ToArray() []T // Alias for ToSlice

			ToSet() map[T]interface{} // Unique elements as set

			// Map terminals with different strategies
			ToMapBy[K comparable](keySelector func(T) K) map[K]T
			ToMapWithValue[K comparable, V any](keySelector func(T) K, valueSelector func(T) V) map[K]V
			ToGroupedMap[K comparable](keySelector func(T) K) map[K][]T
			ToLookup[K comparable, V any](keySelector func(T) K, valueSelector func(T) V) map[K][]V

			// String terminals
			ToString() string
			ToStringWithSeparator(separator string) string
			ToStringWithFormat(format func(T) string) string
			Join(separator string) string

			// Statistical terminals (for numeric types)
			Sum() T
			Average() float64
			Min() (T, error)
			Max() (T, error)
			MinBy[K OrderedStream[K]](keySelector func(T) K) (T, error)
			MaxBy(keySelector func(args interface{}) T) (T, error)

			// Advanced terminals
			ToJSON() (string, error)
			ToCSV() string
			ToCSVWithHeaders(headers []string) string

			// Partitioning terminals
			Partition(predicate func(T) bool) ([]T, []T)
			GroupBy[K comparable](keySelector func(T) K) map[K][]T

			// Existence terminals
			Contains(value T) bool
			ContainsBy(predicate func(T) bool) bool

			// Ordering terminals
			IsSorted() bool
			IsSortedBy(keySelector func(T) K) bool
	*/
}
