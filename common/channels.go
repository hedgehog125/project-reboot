package common

func InitPoolChannel[T any](values ...T) chan T {
	channel := make(chan T, len(values))
	for _, value := range values {
		channel <- value
	}

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
