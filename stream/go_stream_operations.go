package stream

import (
	"gostream/stream/delegates"
	"gostream/stream/operations"
)

func (g *GoStream[T]) Filter(predicate delegates.Predicate[T]) *GoStream[T] {
	g.withOperation(operations.AsFilterOperation(predicate))
	return g
}

func (g *GoStream[T]) Map(mapper func(T) interface{}) *GoStream[T] {
	g.withOperation(operations.AsInterfaceMapOperation(mapper))
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
