package iterators

type (
	Iterator[T any] interface {
		// HasNext returns true if there are more elements
		HasNext() bool

		// Next returns the next element and advances the iterator
		Next() T

		// Close releases any resources held by the iterator
		Close() error
	}

	MapIterator[T, U any] interface {
		Iterator[T]
	}
)
