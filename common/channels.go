package common

func InitChannel[T any](value T) chan T {
	channel := make(chan T)
	go func() { channel <- value }()

	return channel
}
func NewCallbackChannel(callback func()) chan struct{} {
	channel := make(chan struct{})
	go func() {
		callback()
		channel <- struct{}{}
	}()

	return channel
}
