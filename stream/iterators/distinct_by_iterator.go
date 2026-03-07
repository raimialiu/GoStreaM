package iterators

type DistinctByIterator[T, K any] struct {
	source    Iterator[T]
	seen      map[interface{}]bool
	fieldFunc func(T) K
	next      T
	hasNext   bool
	computed  bool
}

func AsDistinctByIterator[T, K any](source Iterator[T], fieldFunc func(T) K) *DistinctByIterator[T, K] {
	return &DistinctByIterator[T, K]{
		source:    source,
		seen:      make(map[interface{}]bool),
		fieldFunc: fieldFunc,
	}
}

func (d *DistinctByIterator[T, K]) HasNext() bool {
	if !d.computed {
		d.computeNext()
	}
	return d.hasNext
}

func (d *DistinctByIterator[T, K]) Next() T {
	if !d.HasNext() {
		var zero T
		return zero
	}

	result := d.next
	d.computed = false
	d.hasNext = false
	return result
}

func (d *DistinctByIterator[T, K]) computeNext() {
	for d.source.HasNext() {
		value := d.source.Next()
		key := d.fieldFunc(value)
		if _, exists := d.seen[key]; !exists {
			d.seen[key] = true
			d.next = value
			d.hasNext = true
			d.computed = true
			return
		}
	}
	d.hasNext = false
	d.computed = true
}

func (d *DistinctByIterator[T, K]) Close() error {
	return d.source.Close()
}
