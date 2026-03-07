package builders

import "github.com/raimialiu/gostream/stream"

type (
	Builder[T any] struct {
		elements []T
		capacity int
	}

	ConditionalBuilder[T any] struct {
		builder   *Builder[T]
		condition bool
		active    bool
	}

	SwitchBuilder[T any] struct {
		builder *Builder[T]
		value   interface{}
		matched bool
	}

	BuilderInterface[T any] interface {
		// Basic addition methods
		Add(element T) *Builder[T]
		AddAll(elements ...T) *Builder[T]
		AddSlice(slice []T) *Builder[T]

		// Conditional addition methods
		AddIf(condition bool, element T) *Builder[T]
		AddIfElse(condition bool, trueElement, falseElement T) *Builder[T]
		AddAllIf(condition bool, elements ...T) *Builder[T]
		AddSliceIf(condition bool, slice []T) *Builder[T]

		// Functional addition methods
		AddWith(supplier func() T) *Builder[T]
		AddWithIf(condition bool, supplier func() T) *Builder[T]
		AddRange(start, count int, generator func(int) T) *Builder[T]
		AddRangeIf(condition bool, start, count int, generator func(int) T) *Builder[T]

		// Conditional logic methods
		When(condition bool) *ConditionalBuilder[T]
		Unless(condition bool) *ConditionalBuilder[T]
		Switch(value interface{}) *SwitchBuilder[T]

		// Stream integration methods
		AddStream(stream stream.GoStream[T]) *Builder[T]
		AddStreamIf(condition bool, stream stream.GoStream[T]) *Builder[T]
		Merge(other *Builder[T]) *Builder[T]

		// Inspection methods
		Count() int
		IsEmpty() bool
		Contains(element T) bool
		Last() (T, bool)

		// Modification methods
		Clear() *Builder[T]
		RemoveLast() *Builder[T]
		RemoveIf(predicate func(T) bool) *Builder[T]

		// Terminal methods
		Build() stream.GoStream[T]
		BuildAndClear() stream.GoStream[T]
		ToSlice() []T
		ToStream() stream.GoStream[T]
	}
)

func NewBuilder[T any]() *Builder[T] {
	return &Builder[T]{
		elements: make([]T, 0),
		capacity: 0,
	}
}

func NewBuilderWithCapacity[T any](capacity int) *Builder[T] {
	return &Builder[T]{
		elements: make([]T, 0, capacity),
		capacity: capacity,
	}
}

// NewBuilderFrom creates a builder from existing elements
func NewBuilderFrom[T any](elements ...T) *Builder[T] {
	return &Builder[T]{
		elements: append([]T{}, elements...),
		capacity: len(elements),
	}
}

func NewBuilderFromSlice[T any](slice []T) *Builder[T] {
	return &Builder[T]{
		elements: append([]T{}, slice...),
		capacity: len(slice),
	}
}

func (b *Builder[T]) Add(element T) *Builder[T] {
	b.elements = append(b.elements, element)
	return b
}

// AddAll adds multiple elements to the builder
func (b *Builder[T]) AddAll(elements ...T) *Builder[T] {
	b.elements = append(b.elements, elements...)
	return b
}

// AddSlice adds all elements from a slice
func (b *Builder[T]) AddSlice(slice []T) *Builder[T] {
	b.elements = append(b.elements, slice...)
	return b
}

// AddFirst adds an element at the beginning
func (b *Builder[T]) AddFirst(element T) *Builder[T] {
	b.elements = append([]T{element}, b.elements...)
	return b
}

// AddAt adds an element at a specific index
func (b *Builder[T]) AddAt(index int, element T) *Builder[T] {
	if index < 0 || index > len(b.elements) {
		return b // Invalid index, ignore
	}

	// Insert at index
	b.elements = append(b.elements[:index], append([]T{element}, b.elements[index:]...)...)
	return b
}

func (b *Builder[T]) AddIf(predicate bool, item T) *Builder[T] {
	if predicate {
		b.elements = append(b.elements, item)
	}

	return b
}

func (b *Builder[T]) AddAllIf(predicate bool, item ...T) *Builder[T] {
	if predicate {
		b.elements = append(b.elements, item...)
	}

	return b
}

func (b *Builder[T]) AddUnless(condition bool, element T) *Builder[T] {
	return b.AddIf(!condition, element)
}

func (b *Builder[T]) AddAllUnless(condition bool, elements ...T) *Builder[T] {
	return b.AddAllIf(!condition, elements...)
}

func (b *Builder[T]) AddWith(supplier func() T) *Builder[T] {
	return b
}
