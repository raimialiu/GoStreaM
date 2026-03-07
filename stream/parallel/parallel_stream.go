package parallel

import (
	"github.com/raimialiu/gostream/stream"
	"runtime"
	"sync"
	"sync/atomic"
)

// ParallelStream wraps a GoStream and executes operations concurrently.
type ParallelStream[T any] struct {
	source  *stream.GoStream[T]
	workers int
}

// From creates a ParallelStream from a GoStream.
func From[T any](s *stream.GoStream[T]) *ParallelStream[T] {
	return &ParallelStream[T]{
		source:  s,
		workers: runtime.NumCPU(),
	}
}

// FromSlice creates a ParallelStream directly from a slice.
func FromSlice[T any](items []T) *ParallelStream[T] {
	return From(stream.From(items))
}

// WithWorkers sets the number of worker goroutines.
func (ps *ParallelStream[T]) WithWorkers(count int) *ParallelStream[T] {
	if count < 1 {
		count = 1
	}
	ps.workers = count
	return ps
}

// Sequential converts back to a regular GoStream.
func (ps *ParallelStream[T]) Sequential() *stream.GoStream[T] {
	return ps.source
}

func (ps *ParallelStream[T]) partition() [][]T {
	items := ps.source.ToSlice()
	if len(items) == 0 {
		return nil
	}
	workers := ps.workers
	if workers > len(items) {
		workers = len(items)
	}
	chunkSize := (len(items) + workers - 1) / workers
	var chunks [][]T
	for i := 0; i < len(items); i += chunkSize {
		end := i + chunkSize
		if end > len(items) {
			end = len(items)
		}
		chunks = append(chunks, items[i:end])
	}
	return chunks
}

// ForEach executes action on each element in parallel.
func (ps *ParallelStream[T]) ForEach(action func(T)) {
	chunks := ps.partition()
	var wg sync.WaitGroup
	for _, chunk := range chunks {
		wg.Add(1)
		go func(items []T) {
			defer wg.Done()
			for _, item := range items {
				action(item)
			}
		}(chunk)
	}
	wg.Wait()
}

// Map applies mapper in parallel and returns results.
func Map[T any, R any](ps *ParallelStream[T], mapper func(T) R) []R {
	items := ps.source.ToSlice()
	if len(items) == 0 {
		return nil
	}
	results := make([]R, len(items))
	workers := ps.workers
	if workers > len(items) {
		workers = len(items)
	}
	chunkSize := (len(items) + workers - 1) / workers

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		start := w * chunkSize
		end := start + chunkSize
		if end > len(items) {
			end = len(items)
		}
		if start >= len(items) {
			break
		}
		wg.Add(1)
		go func(s, e int) {
			defer wg.Done()
			for i := s; i < e; i++ {
				results[i] = mapper(items[i])
			}
		}(start, end)
	}
	wg.Wait()
	return results
}

// Filter filters elements in parallel and returns matching ones.
// Note: order is preserved.
func (ps *ParallelStream[T]) Filter(predicate func(T) bool) *ParallelStream[T] {
	items := ps.source.ToSlice()
	if len(items) == 0 {
		return FromSlice[T](nil)
	}

	workers := ps.workers
	if workers > len(items) {
		workers = len(items)
	}

	matches := make([]bool, len(items))
	chunkSize := (len(items) + workers - 1) / workers

	var wg sync.WaitGroup
	for w := 0; w < workers; w++ {
		start := w * chunkSize
		end := start + chunkSize
		if end > len(items) {
			end = len(items)
		}
		if start >= len(items) {
			break
		}
		wg.Add(1)
		go func(s, e int) {
			defer wg.Done()
			for i := s; i < e; i++ {
				matches[i] = predicate(items[i])
			}
		}(start, end)
	}
	wg.Wait()

	var result []T
	for i, match := range matches {
		if match {
			result = append(result, items[i])
		}
	}
	return FromSlice(result)
}

// AnyMatch returns true if any element matches the predicate (parallel).
func (ps *ParallelStream[T]) AnyMatch(predicate func(T) bool) bool {
	chunks := ps.partition()
	if len(chunks) == 0 {
		return false
	}

	var found atomic.Bool
	var wg sync.WaitGroup
	for _, chunk := range chunks {
		wg.Add(1)
		go func(items []T) {
			defer wg.Done()
			for _, item := range items {
				if found.Load() {
					return
				}
				if predicate(item) {
					found.Store(true)
					return
				}
			}
		}(chunk)
	}
	wg.Wait()
	return found.Load()
}

// AllMatch returns true if all elements match the predicate (parallel).
func (ps *ParallelStream[T]) AllMatch(predicate func(T) bool) bool {
	chunks := ps.partition()
	if len(chunks) == 0 {
		return true
	}

	var failed atomic.Bool
	var wg sync.WaitGroup
	for _, chunk := range chunks {
		wg.Add(1)
		go func(items []T) {
			defer wg.Done()
			for _, item := range items {
				if failed.Load() {
					return
				}
				if !predicate(item) {
					failed.Store(true)
					return
				}
			}
		}(chunk)
	}
	wg.Wait()
	return !failed.Load()
}

// Reduce combines elements in parallel using the accumulator.
func (ps *ParallelStream[T]) Reduce(identity T, accumulator func(T, T) T) T {
	chunks := ps.partition()
	if len(chunks) == 0 {
		return identity
	}

	partials := make([]T, len(chunks))
	var wg sync.WaitGroup
	for i, chunk := range chunks {
		wg.Add(1)
		go func(idx int, items []T) {
			defer wg.Done()
			result := identity
			for _, item := range items {
				result = accumulator(result, item)
			}
			partials[idx] = result
		}(i, chunk)
	}
	wg.Wait()

	result := identity
	for _, partial := range partials {
		result = accumulator(result, partial)
	}
	return result
}

// Count returns the number of elements.
func (ps *ParallelStream[T]) Count() int {
	return len(ps.source.ToSlice())
}

// ToSlice collects all elements into a slice.
func (ps *ParallelStream[T]) ToSlice() []T {
	return ps.source.ToSlice()
}

// Collect applies a collector function across partitions in parallel.
func Collect[T any, R any](ps *ParallelStream[T], collector func([]T) R, combiner func(R, R) R, identity R) R {
	chunks := ps.partition()
	if len(chunks) == 0 {
		return identity
	}

	partials := make([]R, len(chunks))
	var wg sync.WaitGroup
	for i, chunk := range chunks {
		wg.Add(1)
		go func(idx int, items []T) {
			defer wg.Done()
			partials[idx] = collector(items)
		}(i, chunk)
	}
	wg.Wait()

	result := identity
	for _, partial := range partials {
		result = combiner(result, partial)
	}
	return result
}
