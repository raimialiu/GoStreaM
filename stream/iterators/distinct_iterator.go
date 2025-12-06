package iterators

type DistinctIterator[T any] struct {
	source Iterator[T]
	seen   map[interface{}]bool
	index  int
}

func (d DistinctIterator[T]) HasNext() bool {
	return d.source.HasNext()
}

func (d DistinctIterator[T]) Next() T {
	if d.HasNext() {
		// 1, 2, 3, 4, 2, 5, 6, 4, 7, 8 , 9, 10, 1
		value := d.source.Next()
		if _, exist := d.seen[value]; !exist {
			d.seen[value] = true
			return value
		}

		d.index++
		return d.Next()
	}

	var zero T
	return zero
}

func (d DistinctIterator[T]) Close() error {
	d.source.Close()
	return nil
}

func AsDistinctIterator[T any](source Iterator[T]) *DistinctIterator[T] {
	return &DistinctIterator[T]{
		source: source,
		seen:   make(map[interface{}]bool),
		index:  0,
	}
}
