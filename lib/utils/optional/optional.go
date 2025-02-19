package optional

type Optional[T any] struct {
	value    T
	hasValue bool
}

func Of[T any](value T) Optional[T] {
	return Optional[T]{
		value:    value,
		hasValue: true,
	}
}

func WrapString[T ~string](value T) Optional[T] {
	if value == "" {
		return Empty[T]()
	}
	return Of(value)
}

func Empty[T any]() Optional[T] {
	return Optional[T]{
		hasValue: false,
	}
}

func (o Optional[T]) HasValue() bool {
	return o.hasValue
}

func (o Optional[T]) IsEmpty() bool {
	return !o.hasValue
}

func (o Optional[T]) Value() T {
	if !o.hasValue {
		panic("Optional has no value")
	}
	return o.value
}
