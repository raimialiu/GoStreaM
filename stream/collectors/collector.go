package collectors

// Collector defines the interface for collecting stream elements into a result.
// T is the element type, A is the accumulator type, R is the result type.
type Collector[T any, A any, R any] interface {
	// Supplier creates a new accumulator instance
	Supplier() A

	// Accumulator adds an element to the accumulator
	Accumulator(accumulator A, element T) A

	// Finisher transforms the accumulator into the final result
	Finisher(accumulator A) R
}
