package iterators

type IntRageIterator struct {
	start    int
	end      int
	iterator *RangeIterator[int]
}

func AsIntRangeIterator(start, end, step int) *IntRageIterator {
	return &IntRageIterator{
		start:    start,
		end:      end,
		iterator: AsRangeIterator[int](start, end, step),
	}
}

func (it *IntRageIterator) Next() int {
	return it.iterator.Next()
}

func (it *IntRageIterator) HasNext() bool {
	return it.iterator.HasNext()
}

func (it *IntRageIterator) Close() error {
	return it.iterator.Close()
}
