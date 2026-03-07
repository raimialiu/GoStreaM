package iterators

import "github.com/raimialiu/gostream/stream/types"

type RangeIterator[T types.Numeric] struct {
	step  int
	end   int
	item  T
	items []T
	start int
	base  *BaseIterator[T]
}

func AsRangeIterator[T types.Numeric](start, stop, step int) *RangeIterator[T] {
	return &RangeIterator[T]{
		step:  step,
		end:   stop,
		start: start,
		items: make([]T, 0),
		base:  AsBaseIterator[T](),
	}
}

func (it *RangeIterator[T]) Next() T {
	if !it.HasNext() {
		var zero T
		return zero
	}
	value := T(it.start)
	it.items = append(it.items, value)
	it.start += it.step

	return value
}

func (it *RangeIterator[T]) HasNext() bool {
	return it.start < it.end
}

func (it *RangeIterator[T]) Close() error {
	it.start = it.end
	return nil
}
