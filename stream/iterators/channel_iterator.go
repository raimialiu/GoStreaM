package iterators

type ChannelIterator[T any] struct {
	value  <-chan T
	buffer T
	done   bool // to check if channel is not closed
}

func AsChannelIterator[T any](channel <-chan T) *ChannelIterator[T] {
	return &ChannelIterator[T]{
		value: channel,
		done:  false,
	}
}

func (it *ChannelIterator[T]) Next() T {
	if !it.HasNext() {
		var zero T
		return zero
	}

	value := it.buffer
	it.readFromChannel()

	return value
}

func (it *ChannelIterator[T]) readFromChannel() T {
	if it.done {
		var zero T
		return zero
	}

	value, ok := <-it.value
	if !ok {
		it.done = true
	}

	return value
}

func (it *ChannelIterator[T]) HasNext() bool {
	return !it.done
}

func (it *ChannelIterator[T]) Close() error {
	it.done = true
	return nil
}
