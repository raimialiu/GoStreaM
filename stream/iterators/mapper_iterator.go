package iterators

import "reflect"

// MapperIterator transforms elements from type T to type R
type MapperIterator[T, R any] struct {
	mapper  func(T) R
	source  Iterator[T]
	next    R
	hasNext bool
}

func AsMapperIterator[T, R any](mapper func(T) R, source Iterator[T]) *MapperIterator[T, R] {
	return &MapperIterator[T, R]{
		mapper:  mapper,
		hasNext: false,
		source:  source,
	}
}

func (it *MapperIterator[T, R]) HasNext() bool {
	if !it.hasNext {
		it.mapNext()
	}
	return it.hasNext
}

func (it *MapperIterator[T, R]) Next() R {
	if !it.HasNext() {
		var zero R
		return zero
	}

	current := it.next
	it.hasNext = false // Reset for next iteration
	return current
}

func (it *MapperIterator[T, R]) mapNext() {
	var zero R
	for it.source.HasNext() {
		sourceItem := it.source.Next()
		mapperResult := it.mapper(sourceItem)

		// Check if the mapped result is not zero value
		if !reflect.DeepEqual(mapperResult, zero) {
			it.next = mapperResult
			it.hasNext = true
			return
		}
	}
	it.hasNext = false
}

func (it *MapperIterator[T, R]) Close() error {
	return it.source.Close()
}
