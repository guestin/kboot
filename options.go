package kboot

type Option[T any] interface {
	apply(t T)
}

type optionFunc[T any] func(t T)

func (f optionFunc[T]) apply(t T) {
	f(t)
}
