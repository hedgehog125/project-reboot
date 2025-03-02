package util

type ErrWithIndex struct {
	Err   error
	Index int
}

type ErrWithPointer[T any] struct {
	Err     error
	Pointer T
}
