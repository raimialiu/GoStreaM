package iterators

type DistinctIterator[T any] struct {
	source   Iterator[T]
	seen     map[interface{}]bool
	next     T
	hasNext  bool
	computed bool
}

func AsDistinctIterator[T any](source Iterator[T]) *DistinctIterator[T] {
	return &DistinctIterator[T]{
		source: source,
		seen:   make(map[interface{}]bool),
	}
}

func (d *DistinctIterator[T]) HasNext() bool {
	if !d.computed {
		d.computeNext()
	}
	return d.hasNext
}

func (d *DistinctIterator[T]) Next() T {
	if !d.HasNext() {
		var zero T
		return zero
	}

	result := d.next
	d.computed = false
	d.hasNext = false
	return result
}

func (d *DistinctIterator[T]) computeNext() {
	for d.source.HasNext() {
		value := d.source.Next()
		if _, exists := d.seen[value]; !exists {
			d.seen[value] = true
			d.next = value
			d.hasNext = true
			d.computed = true
			return
		}
	}
	d.hasNext = false
	d.computed = true
}

func (d *DistinctIterator[T]) Close() error {
	return d.source.Close()
}
