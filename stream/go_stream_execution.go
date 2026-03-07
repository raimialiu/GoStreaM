package stream

import (
	"github.com/raimialiu/gostream/stream/iterators"
	"reflect"
	"sort"
)

func (g *GoStream[T]) execute() []iterators.Iterator[T] {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.closed {
		panic("cannot execute closed stream")
	}

	iters := make([]iterators.Iterator[T], 0, len(g.iterators))
	for _, iterator := range g.iterators {
		iter := iterator
		for _, op := range g.operations {
			iter = op.Apply(iter)
		}
		iters = append(iters, iter)
	}

	return iters
}

func (g *GoStream[T]) collect() []T {
	var list []T
	for _, source := range g.execute() {
		for source.HasNext() {
			list = append(list, source.Next())
		}
	}
	return list
}

// ========== TERMINAL OPERATIONS ==========

func (g *GoStream[T]) Count() int {
	count := 0
	for _, source := range g.execute() {
		for source.HasNext() {
			source.Next()
			count++
		}
	}
	return count
}

func (g *GoStream[T]) CountBy(predicate func(T) bool) int {
	count := 0
	for _, source := range g.execute() {
		for source.HasNext() {
			if predicate(source.Next()) {
				count++
			}
		}
	}
	return count
}

func (g *GoStream[T]) Sum() float64 {
	sum := 0.0
	for _, source := range g.execute() {
		for source.HasNext() {
			value := source.Next()
			rv := reflect.ValueOf(value)
			switch rv.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				sum += float64(rv.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				sum += float64(rv.Uint())
			case reflect.Float32, reflect.Float64:
				sum += rv.Float()
			}
		}
	}
	return sum
}

func (g *GoStream[T]) Average() float64 {
	sum := 0.0
	count := 0
	for _, source := range g.execute() {
		for source.HasNext() {
			value := source.Next()
			rv := reflect.ValueOf(value)
			switch rv.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				sum += float64(rv.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				sum += float64(rv.Uint())
			case reflect.Float32, reflect.Float64:
				sum += rv.Float()
			}
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return sum / float64(count)
}

func (g *GoStream[T]) Min() (T, error) {
	var min T
	found := false
	for _, source := range g.execute() {
		for source.HasNext() {
			value := source.Next()
			if !found {
				min = value
				found = true
				continue
			}
			rv := reflect.ValueOf(value)
			rmv := reflect.ValueOf(min)
			switch rv.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if rv.Int() < rmv.Int() {
					min = value
				}
			case reflect.Float32, reflect.Float64:
				if rv.Float() < rmv.Float() {
					min = value
				}
			case reflect.String:
				if rv.String() < rmv.String() {
					min = value
				}
			}
		}
	}
	if !found {
		return min, ErrEmptyStream
	}
	return min, nil
}

func (g *GoStream[T]) Max() (T, error) {
	var max T
	found := false
	for _, source := range g.execute() {
		for source.HasNext() {
			value := source.Next()
			if !found {
				max = value
				found = true
				continue
			}
			rv := reflect.ValueOf(value)
			rmv := reflect.ValueOf(max)
			switch rv.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if rv.Int() > rmv.Int() {
					max = value
				}
			case reflect.Float32, reflect.Float64:
				if rv.Float() > rmv.Float() {
					max = value
				}
			case reflect.String:
				if rv.String() > rmv.String() {
					max = value
				}
			}
		}
	}
	if !found {
		return max, ErrEmptyStream
	}
	return max, nil
}

func (g *GoStream[T]) ToList() []T {
	return g.collect()
}

func (g *GoStream[T]) ToSlice() []T {
	return g.collect()
}

func (g *GoStream[T]) ToMap(keySelector func(T) interface{}) map[interface{}]T {
	result := make(map[interface{}]T)
	for _, source := range g.execute() {
		for source.HasNext() {
			value := source.Next()
			key := keySelector(value)
			result[key] = value
		}
	}
	return result
}

func (g *GoStream[T]) GroupBy(keySelector func(T) interface{}) map[interface{}][]T {
	result := make(map[interface{}][]T)
	for _, source := range g.execute() {
		for source.HasNext() {
			value := source.Next()
			key := keySelector(value)
			result[key] = append(result[key], value)
		}
	}
	return result
}

func (g *GoStream[T]) Partition(predicate func(T) bool) ([]T, []T) {
	var matching, nonMatching []T
	for _, source := range g.execute() {
		for source.HasNext() {
			value := source.Next()
			if predicate(value) {
				matching = append(matching, value)
			} else {
				nonMatching = append(nonMatching, value)
			}
		}
	}
	return matching, nonMatching
}

func (g *GoStream[T]) ForEach(action func(T)) {
	for _, source := range g.execute() {
		for source.HasNext() {
			action(source.Next())
		}
	}
}

func (g *GoStream[T]) Any() bool {
	for _, source := range g.execute() {
		if source.HasNext() {
			return true
		}
	}
	return false
}

func (g *GoStream[T]) AnyMatch(predicate func(T) bool) bool {
	for _, source := range g.execute() {
		for source.HasNext() {
			if predicate(source.Next()) {
				return true
			}
		}
	}
	return false
}

func (g *GoStream[T]) AllMatch(predicate func(T) bool) bool {
	for _, source := range g.execute() {
		for source.HasNext() {
			if !predicate(source.Next()) {
				return false
			}
		}
	}
	return true
}

func (g *GoStream[T]) NoneMatch(predicate func(T) bool) bool {
	return !g.AnyMatch(predicate)
}

func (g *GoStream[T]) First() (T, error) {
	for _, source := range g.execute() {
		if source.HasNext() {
			return source.Next(), nil
		}
	}
	var zero T
	return zero, ErrEmptyStream
}

func (g *GoStream[T]) FirstOrDefault(defaultValue T) T {
	val, err := g.First()
	if err != nil {
		return defaultValue
	}
	return val
}

func (g *GoStream[T]) Last() (T, error) {
	var last T
	found := false
	for _, source := range g.execute() {
		for source.HasNext() {
			last = source.Next()
			found = true
		}
	}
	if !found {
		return last, ErrEmptyStream
	}
	return last, nil
}

func (g *GoStream[T]) Single() (T, error) {
	var result T
	count := 0
	for _, source := range g.execute() {
		for source.HasNext() {
			result = source.Next()
			count++
			if count > 1 {
				var zero T
				return zero, ErrMultipleElements
			}
		}
	}
	if count == 0 {
		return result, ErrEmptyStream
	}
	return result, nil
}

func (g *GoStream[T]) Reduce(identity T, accumulator func(T, T) T) T {
	result := identity
	for _, source := range g.execute() {
		for source.HasNext() {
			result = accumulator(result, source.Next())
		}
	}
	return result
}

func (g *GoStream[T]) Contains(value T) bool {
	for _, source := range g.execute() {
		for source.HasNext() {
			if reflect.DeepEqual(source.Next(), value) {
				return true
			}
		}
	}
	return false
}

func (g *GoStream[T]) OrderBy(less func(a, b T) bool) *GoStream[T] {
	items := g.collect()
	sort.Slice(items, func(i, j int) bool {
		return less(items[i], items[j])
	})
	return From(items)
}

func (g *GoStream[T]) OrderByDesc(less func(a, b T) bool) *GoStream[T] {
	items := g.collect()
	sort.Slice(items, func(i, j int) bool {
		return less(items[j], items[i])
	})
	return From(items)
}

func (g *GoStream[T]) Close() {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.closed {
		return
	}
	g.closed = true
}
