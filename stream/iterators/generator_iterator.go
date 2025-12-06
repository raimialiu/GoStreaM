package iterators

type GeneratorIterator[T any] struct {
	valueFunc    func(args ...interface{}) T
	currentValue T
	valueArgs    []interface{}
}

func AsGeneratorIterator[T any](valueFunc func(args ...interface{}) T) *GeneratorIterator[T] {
	var zero T
	return &GeneratorIterator[T]{
		valueFunc:    valueFunc,
		currentValue: zero,
		valueArgs:    make([]interface{}, 0),
	}
}

// Will directly execute the function logic
func (it *GeneratorIterator[T]) WithNextArgs(args ...interface{}) *GeneratorIterator[T] {
	it.currentValue = it.valueFunc(args...)
	return it
}

// Will not execut the function logic, but setup the arguments till the next time you call Next
func (it *GeneratorIterator[T]) WithArgs(args ...interface{}) *GeneratorIterator[T] {
	it.valueArgs = append(it.valueArgs, args...)
	return it
}

func (it *GeneratorIterator[T]) Next() T {
	if !it.HasNext() {
		var zero T
		return zero
	}

	/*
		if it.valueFunc != nil {
			return it.currentValue
		}
	*/

	return it.valueFunc(it.valueArgs...)
}

func (it *GeneratorIterator[T]) HasNext() bool {
	return it.valueFunc != nil
}

func (it *GeneratorIterator[T]) Close() error {
	it.valueFunc = nil
	return nil
}
