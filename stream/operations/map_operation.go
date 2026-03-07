package operations

import "github.com/raimialiu/gostream/stream/iterators"

// Generic map operation that can transform types
type GenericMapOperation[T, R any] struct {
	mapper func(T) R
}

func AsGenericMapOperation[T, R any](mapper func(T) R) *GenericMapOperation[T, R] {
	return &GenericMapOperation[T, R]{
		mapper: mapper,
	}
}

func (f GenericMapOperation[T, R]) Apply(source iterators.Iterator[T]) iterators.Iterator[R] {
	return iterators.AsMapperIterator[T, R](f.mapper, source)
}

// Specific map operation for interface{} that maintains type T
type InterfaceMapOperation[T any] struct {
	mapper func(T) interface{}
}

func AsInterfaceMapOperation[T any](mapper func(T) interface{}) *InterfaceMapOperation[T] {
	return &InterfaceMapOperation[T]{
		mapper: mapper,
	}
}

func (f InterfaceMapOperation[T]) CanFuse(next StreamOperation[T]) bool {
	return true
}

func (f InterfaceMapOperation[T]) Apply(source iterators.Iterator[T]) iterators.Iterator[T] {
	// Use the type asserting iterator from Solution 1
	return &TypeAssertingIterator[T]{
		source: iterators.AsMapperIterator[T, interface{}](f.mapper, source),
	}
}

func (f InterfaceMapOperation[T]) ApplyFunc(func(args ...interface{}) interface{}) iterators.Iterator[T] {
	//TODO implement me
	panic("implement me")
}

// Helper iterator that performs type assertion
type TypeAssertingIterator[T any] struct {
	source iterators.Iterator[interface{}]
}

func (it *TypeAssertingIterator[T]) HasNext() bool {
	return it.source.HasNext()
}

func (it *TypeAssertingIterator[T]) Next() T {
	value := it.source.Next()
	if typedValue, ok := value.(T); ok {
		return typedValue
	}
	var zero T
	return zero
}

func (it *TypeAssertingIterator[T]) Close() error {
	return it.source.Close()
}
