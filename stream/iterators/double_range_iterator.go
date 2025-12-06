package iterators

type DoubleRangeIterator struct {
	start    float32
	end      float32
	iterator *RangeIterator[float32]
}

func AsDoubleRangeIterator(start, end, step float32) *DoubleRangeIterator {
	begin := int(start)
	stop := int(end)
	return &DoubleRangeIterator{
		start:    start,
		end:      end,
		iterator: AsRangeIterator[float32](begin, stop, int(step)),
	}
}

func (it *DoubleRangeIterator) Next() float32 {
	return it.iterator.Next()
}

func (it *DoubleRangeIterator) HasNext() bool {
	return it.iterator.HasNext()
}

func (it *DoubleRangeIterator) Close() error {
	return it.iterator.Close()
}
