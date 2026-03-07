package operations

import "github.com/raimialiu/gostream/stream/iterators"

// StreamOperation represents a lazy intermediate operation in the stream pipeline.
// Implement this interface to create custom operations.
type StreamOperation[T any] interface {
	// CanFuse returns true if this operation can be fused with the next
	CanFuse(next StreamOperation[T]) bool

	// Apply wraps the source iterator with this operation's logic
	Apply(source iterators.Iterator[T]) iterators.Iterator[T]
}
