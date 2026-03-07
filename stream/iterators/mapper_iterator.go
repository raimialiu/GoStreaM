package iterators

type MapperIterator[T, R any] struct {
	mapper   func(T) R
	source   Iterator[T]
	next     R
	hasNext  bool
	computed bool
}

func AsMapperIterator[T, R any](mapper func(T) R, source Iterator[T]) *MapperIterator[T, R] {
	return &MapperIterator[T, R]{
		mapper: mapper,
		source: source,
	}
}

func (it *MapperIterator[T, R]) HasNext() bool {
	if !it.computed {
		it.computeNext()
	}
	return it.hasNext
}

func (it *MapperIterator[T, R]) Next() R {
	if !it.HasNext() {
		var zero R
		return zero
	}

	result := it.next
	it.computed = false
	it.hasNext = false
	return result
}

func (it *MapperIterator[T, R]) computeNext() {
	if it.source.HasNext() {
		sourceItem := it.source.Next()
		it.next = it.mapper(sourceItem)
		it.hasNext = true
	} else {
		it.hasNext = false
	}
	it.computed = true
}

func (it *MapperIterator[T, R]) Close() error {
	return it.source.Close()
}
