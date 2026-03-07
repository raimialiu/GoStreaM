package iterators

type FlatMapIterator[T any] struct {
	source   Iterator[T]
	mapper   func(T) []T
	current  Iterator[T]
	hasNext  bool
	computed bool
	next     T
}

func AsFlatMapIterator[T any](source Iterator[T], mapper func(T) []T) *FlatMapIterator[T] {
	return &FlatMapIterator[T]{
		source: source,
		mapper: mapper,
	}
}

func (it *FlatMapIterator[T]) HasNext() bool {
	if !it.computed {
		it.computeNext()
	}
	return it.hasNext
}

func (it *FlatMapIterator[T]) Next() T {
	if !it.HasNext() {
		var zero T
		return zero
	}
	result := it.next
	it.computed = false
	it.hasNext = false
	return result
}

func (it *FlatMapIterator[T]) computeNext() {
	for {
		if it.current != nil && it.current.HasNext() {
			it.next = it.current.Next()
			it.hasNext = true
			it.computed = true
			return
		}
		if !it.source.HasNext() {
			it.hasNext = false
			it.computed = true
			return
		}
		items := it.mapper(it.source.Next())
		it.current = AsListIterator(items...)
	}
}

func (it *FlatMapIterator[T]) Close() error {
	return it.source.Close()
}
