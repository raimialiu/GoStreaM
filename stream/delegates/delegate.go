package delegates

type (
	Predicate[T any]          func(T) bool
	KeySelector[K any]        func() K
	FuncKeySelector[K, V any] func() K
)
