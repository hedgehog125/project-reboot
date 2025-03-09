package common

// Service types should go in services.go

type ErrWithIndex struct {
	Err   error
	Index int
}

type ErrWithPointer[T any] struct {
	Err     error
	Pointer T
}
