package xatago

type Record interface {
	ID() string
}

func SetPtr[T any](value T) *T {
	t := new(T)
	*t = value
	return t
}
