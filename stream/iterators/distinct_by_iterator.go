package iterators

type DistinctByIterator[T, K any] struct {
	source    Iterator[T]
	seen      map[interface{}]bool
	index     int
	fieldFunc func(T) K
}

func AsDistinctByIterator[T, K any](source Iterator[T], fieldFunc func(T) K) *DistinctByIterator[T, K] {
	return &DistinctByIterator[T, K]{
		source:    source,
		seen:      make(map[interface{}]bool),
		fieldFunc: fieldFunc,
	}
}

func (d DistinctByIterator[T, K]) HasNext() bool {
	return d.source.HasNext()
}

func (d DistinctByIterator[T, K]) Next() T {
	var zero T
	if !d.HasNext() {
		return zero
	}

	value := d.source.Next()
	fieldFunc := d.fieldFunc(value)

	if _, exists := d.seen[fieldFunc]; !exists {
		d.seen[fieldFunc] = true
		return value
	}

	return d.Next()

}

func (d DistinctByIterator[T, K]) Close() error {
	return d.source.Close()
}
