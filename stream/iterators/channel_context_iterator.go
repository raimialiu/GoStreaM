package iterators

import "context"

type ChannelContextIterator[T any] struct {
	value   <-chan T
	buffer  T
	context context.Context
	done    bool
}

func AsChannelContextIterator[T any](ctx context.Context, value <-chan T) *ChannelContextIterator[T] {
	return &ChannelContextIterator[T]{
		value:   value,
		context: ctx,
	}
}

func (it *ChannelContextIterator[T]) Close() error {
	it.done = true
	return nil
}

func (it *ChannelContextIterator[T]) HasNext() bool {
	return !it.done
}

func (it *ChannelContextIterator[T]) Next() T {
	var zero T
	if !it.HasNext() {
		return zero
	}

	value := it.buffer
	it.readFromChannel()
	return value
}

func (it *ChannelContextIterator[T]) readFromChannel() T {
	var zero T
	if it.done {
		return zero
	}

	for {
		select {
		case <-it.context.Done():
			it.done = true
			return zero
		case item, ok := <-it.value:
			if !ok {
				it.done = true
				return zero
			}
			it.buffer = item
			return it.buffer
		}
	}
}
