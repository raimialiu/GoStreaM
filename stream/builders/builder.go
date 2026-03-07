package builders

import (
	"github.com/raimialiu/gostream/stream"
	"reflect"
)

type Builder[T any] struct {
	elements []T
}

func NewBuilder[T any]() *Builder[T] {
	return &Builder[T]{elements: make([]T, 0)}
}

func NewBuilderWithCapacity[T any](capacity int) *Builder[T] {
	return &Builder[T]{elements: make([]T, 0, capacity)}
}

func NewBuilderFrom[T any](elements ...T) *Builder[T] {
	return &Builder[T]{elements: append([]T{}, elements...)}
}

func NewBuilderFromSlice[T any](slice []T) *Builder[T] {
	return &Builder[T]{elements: append([]T{}, slice...)}
}

// ========== Basic Addition ==========

func (b *Builder[T]) Add(element T) *Builder[T] {
	b.elements = append(b.elements, element)
	return b
}

func (b *Builder[T]) AddAll(elements ...T) *Builder[T] {
	b.elements = append(b.elements, elements...)
	return b
}

func (b *Builder[T]) AddSlice(slice []T) *Builder[T] {
	b.elements = append(b.elements, slice...)
	return b
}

func (b *Builder[T]) AddFirst(element T) *Builder[T] {
	b.elements = append([]T{element}, b.elements...)
	return b
}

func (b *Builder[T]) AddAt(index int, element T) *Builder[T] {
	if index < 0 || index > len(b.elements) {
		return b
	}
	b.elements = append(b.elements[:index], append([]T{element}, b.elements[index:]...)...)
	return b
}

// ========== Conditional Addition ==========

func (b *Builder[T]) AddIf(condition bool, element T) *Builder[T] {
	if condition {
		b.elements = append(b.elements, element)
	}
	return b
}

func (b *Builder[T]) AddIfElse(condition bool, trueElement, falseElement T) *Builder[T] {
	if condition {
		b.elements = append(b.elements, trueElement)
	} else {
		b.elements = append(b.elements, falseElement)
	}
	return b
}

func (b *Builder[T]) AddAllIf(condition bool, elements ...T) *Builder[T] {
	if condition {
		b.elements = append(b.elements, elements...)
	}
	return b
}

func (b *Builder[T]) AddSliceIf(condition bool, slice []T) *Builder[T] {
	if condition {
		b.elements = append(b.elements, slice...)
	}
	return b
}

func (b *Builder[T]) AddUnless(condition bool, element T) *Builder[T] {
	return b.AddIf(!condition, element)
}

func (b *Builder[T]) AddAllUnless(condition bool, elements ...T) *Builder[T] {
	return b.AddAllIf(!condition, elements...)
}

// ========== Functional Addition ==========

func (b *Builder[T]) AddWith(supplier func() T) *Builder[T] {
	b.elements = append(b.elements, supplier())
	return b
}

func (b *Builder[T]) AddWithIf(condition bool, supplier func() T) *Builder[T] {
	if condition {
		b.elements = append(b.elements, supplier())
	}
	return b
}

func (b *Builder[T]) AddRange(start, count int, generator func(int) T) *Builder[T] {
	for i := start; i < start+count; i++ {
		b.elements = append(b.elements, generator(i))
	}
	return b
}

func (b *Builder[T]) AddRangeIf(condition bool, start, count int, generator func(int) T) *Builder[T] {
	if condition {
		return b.AddRange(start, count, generator)
	}
	return b
}

// ========== Stream Integration ==========

func (b *Builder[T]) AddStream(s *stream.GoStream[T]) *Builder[T] {
	items := s.ToList()
	b.elements = append(b.elements, items...)
	return b
}

func (b *Builder[T]) AddStreamIf(condition bool, s *stream.GoStream[T]) *Builder[T] {
	if condition {
		return b.AddStream(s)
	}
	return b
}

func (b *Builder[T]) Merge(other *Builder[T]) *Builder[T] {
	b.elements = append(b.elements, other.elements...)
	return b
}

// ========== Conditional Logic ==========

func (b *Builder[T]) When(condition bool) *ConditionalBuilder[T] {
	return &ConditionalBuilder[T]{
		builder:   b,
		condition: condition,
	}
}

func (b *Builder[T]) Unless(condition bool) *ConditionalBuilder[T] {
	return b.When(!condition)
}

func (b *Builder[T]) Switch(value interface{}) *SwitchBuilder[T] {
	return &SwitchBuilder[T]{
		builder: b,
		value:   value,
	}
}

// ========== Inspection ==========

func (b *Builder[T]) Count() int {
	return len(b.elements)
}

func (b *Builder[T]) IsEmpty() bool {
	return len(b.elements) == 0
}

func (b *Builder[T]) Contains(element T) bool {
	for _, e := range b.elements {
		if reflect.DeepEqual(e, element) {
			return true
		}
	}
	return false
}

func (b *Builder[T]) Last() (T, bool) {
	if len(b.elements) == 0 {
		var zero T
		return zero, false
	}
	return b.elements[len(b.elements)-1], true
}

func (b *Builder[T]) Get(index int) (T, bool) {
	if index < 0 || index >= len(b.elements) {
		var zero T
		return zero, false
	}
	return b.elements[index], true
}

// ========== Modification ==========

func (b *Builder[T]) Clear() *Builder[T] {
	b.elements = make([]T, 0)
	return b
}

func (b *Builder[T]) RemoveLast() *Builder[T] {
	if len(b.elements) > 0 {
		b.elements = b.elements[:len(b.elements)-1]
	}
	return b
}

func (b *Builder[T]) RemoveAt(index int) *Builder[T] {
	if index < 0 || index >= len(b.elements) {
		return b
	}
	b.elements = append(b.elements[:index], b.elements[index+1:]...)
	return b
}

func (b *Builder[T]) RemoveIf(predicate func(T) bool) *Builder[T] {
	filtered := make([]T, 0, len(b.elements))
	for _, e := range b.elements {
		if !predicate(e) {
			filtered = append(filtered, e)
		}
	}
	b.elements = filtered
	return b
}

// ========== Terminal Methods ==========

func (b *Builder[T]) Build() *stream.GoStream[T] {
	return stream.From(append([]T{}, b.elements...))
}

func (b *Builder[T]) BuildAndClear() *stream.GoStream[T] {
	s := b.Build()
	b.Clear()
	return s
}

func (b *Builder[T]) ToSlice() []T {
	return append([]T{}, b.elements...)
}

func (b *Builder[T]) ToStream() *stream.GoStream[T] {
	return b.Build()
}

// ========== ConditionalBuilder ==========

type ConditionalBuilder[T any] struct {
	builder   *Builder[T]
	condition bool
}

func (cb *ConditionalBuilder[T]) Add(element T) *ConditionalBuilder[T] {
	if cb.condition {
		cb.builder.Add(element)
	}
	return cb
}

func (cb *ConditionalBuilder[T]) AddAll(elements ...T) *ConditionalBuilder[T] {
	if cb.condition {
		cb.builder.AddAll(elements...)
	}
	return cb
}

func (cb *ConditionalBuilder[T]) AddSlice(slice []T) *ConditionalBuilder[T] {
	if cb.condition {
		cb.builder.AddSlice(slice)
	}
	return cb
}

func (cb *ConditionalBuilder[T]) AddWith(supplier func() T) *ConditionalBuilder[T] {
	if cb.condition {
		cb.builder.AddWith(supplier)
	}
	return cb
}

func (cb *ConditionalBuilder[T]) ElseAdd(element T) *ConditionalBuilder[T] {
	if !cb.condition {
		cb.builder.Add(element)
	}
	return cb
}

func (cb *ConditionalBuilder[T]) ElseAddAll(elements ...T) *ConditionalBuilder[T] {
	if !cb.condition {
		cb.builder.AddAll(elements...)
	}
	return cb
}

func (cb *ConditionalBuilder[T]) End() *Builder[T] {
	return cb.builder
}

// ========== SwitchBuilder ==========

type SwitchBuilder[T any] struct {
	builder *Builder[T]
	value   interface{}
	matched bool
}

func (sb *SwitchBuilder[T]) Case(caseValue interface{}, elements ...T) *SwitchBuilder[T] {
	if !sb.matched && reflect.DeepEqual(sb.value, caseValue) {
		sb.builder.AddAll(elements...)
		sb.matched = true
	}
	return sb
}

func (sb *SwitchBuilder[T]) CaseWith(caseValue interface{}, supplier func() T) *SwitchBuilder[T] {
	if !sb.matched && reflect.DeepEqual(sb.value, caseValue) {
		sb.builder.AddWith(supplier)
		sb.matched = true
	}
	return sb
}

func (sb *SwitchBuilder[T]) Default(elements ...T) *SwitchBuilder[T] {
	if !sb.matched {
		sb.builder.AddAll(elements...)
		sb.matched = true
	}
	return sb
}

func (sb *SwitchBuilder[T]) DefaultWith(supplier func() T) *SwitchBuilder[T] {
	if !sb.matched {
		sb.builder.AddWith(supplier)
		sb.matched = true
	}
	return sb
}

func (sb *SwitchBuilder[T]) End() *Builder[T] {
	return sb.builder
}
