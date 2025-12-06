package stream

import (
	"gostream/stream/iterators"
	"reflect"
)

func (g *GoStream[T]) execute() iterators.Iterator[T] {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.closed {
		panic("cannot execute closed stream")
	}

	var iter = g.iterator
	for _, op := range g.operations {
		iter = op.Apply(iter)
	}

	return iter.(iterators.Iterator[T])
}

func (g *GoStream[T]) Count() int {
	var counter int = 0

	var source = g.execute()
	for source.HasNext() {
		counter++
	}

	return counter
}

func (g *GoStream[T]) CountBy(predicate func(T) bool) int {
	counter := 0
	var source = g.execute()
	for source.HasNext() {
		var value = source.Next()
		if predicate(value) {
			counter++
		}
	}

	return counter
}

func (g *GoStream[T]) ToList() []T {
	var list = make([]T, 0)
	var source = g.execute()
	for source.HasNext() {
		value := source.Next()
		var zero T
		if reflect.DeepEqual(value, zero) {
			continue
		}
		list = append(list, value)
	}

	return list
}

func (g *GoStream[T]) Close() {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.closed {
		return
	}

	g.closed = true
}
