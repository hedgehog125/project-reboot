package util

func InitChannel[T any](value T) chan T {
	channel := make(chan T)
	go func() { channel <- value }()

	return channel
}
