package stream

import (
	"github.com/raimialiu/gostream/stream/delegates"
	"github.com/raimialiu/gostream/stream/iterators"
	"github.com/raimialiu/gostream/stream/operations"
)

// ========== INTERMEDIATE OPERATIONS (Lazy) ==========

func (g *GoStream[T]) Filter(predicate delegates.Predicate[T]) *GoStream[T] {
	g.withOperation(operations.AsFilterOperation(predicate))
	return g
}

func (g *GoStream[T]) Map(mapper func(T) interface{}) *GoStream[T] {
	g.withOperation(operations.AsInterfaceMapOperation(mapper))
	return g
}

func (g *GoStream[T]) FlatMap(mapper func(T) []T) *GoStream[T] {
	g.withOperation(operations.AsFlatMapOperation(mapper))
	return g
}

func (g *GoStream[T]) Distinct() *GoStream[T] {
	g.withOperation(operations.AsDistinct[T]())
	return g
}

func (g *GoStream[T]) DistinctBy(fieldFunc func(T) interface{}) *GoStream[T] {
	g.withOperation(operations.AsDistinctBy(fieldFunc))
	return g
}

func (g *GoStream[T]) Select(selector func(T) interface{}) *GoStream[T] {
	return g.Map(selector)
}

func (g *GoStream[T]) Take(count int) *GoStream[T] {
	g.withOperation(operations.AsTakeOperation[T](count))
	return g
}

func (g *GoStream[T]) Skip(count int) *GoStream[T] {
	g.withOperation(operations.AsSkipOperation[T](count))
	return g
}

func (g *GoStream[T]) TakeWhile(predicate func(T) bool) *GoStream[T] {
	g.withOperation(operations.AsTakeWhileOperation(predicate))
	return g
}

func (g *GoStream[T]) SkipWhile(predicate func(T) bool) *GoStream[T] {
	g.withOperation(operations.AsSkipWhileOperation(predicate))
	return g
}

func (g *GoStream[T]) Peek(action func(T)) *GoStream[T] {
	g.withOperation(operations.AsPeekOperation(action))
	return g
}

func (g *GoStream[T]) Reverse() *GoStream[T] {
	items := g.collect()
	for i, j := 0, len(items)-1; i < j; i, j = i+1, j-1 {
		items[i], items[j] = items[j], items[i]
	}
	result := From(items)
	return result
}

func (g *GoStream[T]) Concat(others ...GoStream[T]) *GoStream[T] {
	allIters := make([]iterators.Iterator[T], 0)
	allIters = append(allIters, g.iterators...)
	for _, other := range others {
		allIters = append(allIters, other.iterators...)
	}
	g.iterators = []iterators.Iterator[T]{iterators.AsConcatIterator(allIters...)}
	return g
}

func (g *GoStream[T]) Union(other GoStream[T]) *GoStream[T] {
	return g.Concat(other).Distinct()
}

func (g *GoStream[T]) Intersect(other GoStream[T]) *GoStream[T] {
	otherSet := make(map[interface{}]bool)
	for _, iter := range other.iterators {
		for iter.HasNext() {
			otherSet[iter.Next()] = true
		}
	}
	g.withOperation(operations.AsFilterOperation[T](func(item T) bool {
		_, exists := otherSet[item]
		return exists
	}))
	return g.Distinct()
}

func (g *GoStream[T]) Except(other GoStream[T]) *GoStream[T] {
	otherSet := make(map[interface{}]bool)
	for _, iter := range other.iterators {
		for iter.HasNext() {
			otherSet[iter.Next()] = true
		}
	}
	g.withOperation(operations.AsFilterOperation[T](func(item T) bool {
		_, exists := otherSet[item]
		return !exists
	}))
	return g
}

func (g *GoStream[T]) Cache() *GoStream[T] {
	items := g.collect()
	return From(items)
}

func (g *GoStream[T]) Zip(other *GoStream[T], zipper func(T, T) T) *GoStream[T] {
	items1 := g.collect()
	items2 := other.collect()
	minLen := len(items1)
	if len(items2) < minLen {
		minLen = len(items2)
	}
	result := make([]T, minLen)
	for i := 0; i < minLen; i++ {
		result[i] = zipper(items1[i], items2[i])
	}
	return From(result)
}

func (g *GoStream[T]) Iterator() iterators.Iterator[T] {
	iters := g.execute()
	if len(iters) == 0 {
		return iterators.AsEmptyIterator[T]()
	}
	if len(iters) == 1 {
		return iters[0]
	}
	var all []T
	for _, iter := range iters {
		for iter.HasNext() {
			all = append(all, iter.Next())
		}
	}
	return iterators.AsListIterator(all...)
}

func (g *GoStream[T]) IsParallel() bool {
	return g.parallel
}
