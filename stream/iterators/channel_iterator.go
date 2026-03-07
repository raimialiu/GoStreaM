package iterators

type ChannelIterator[T any] struct {
	ch       <-chan T
	buffer   T
	hasNext  bool
	computed bool
}

func AsChannelIterator[T any](channel <-chan T) *ChannelIterator[T] {
	return &ChannelIterator[T]{
		ch: channel,
	}
}

func (it *ChannelIterator[T]) HasNext() bool {
	if !it.computed {
		it.readNext()
	}
	return it.hasNext
}

func (it *ChannelIterator[T]) Next() T {
	if !it.HasNext() {
		var zero T
		return zero
	}

	value := it.buffer
	it.computed = false
	it.hasNext = false
	return value
}

func (it *ChannelIterator[T]) readNext() {
	value, ok := <-it.ch
	if ok {
		it.buffer = value
		it.hasNext = true
	} else {
		it.hasNext = false
	}
	it.computed = true
}

func (it *ChannelIterator[T]) Close() error {
	it.hasNext = false
	it.computed = true
	return nil
}
