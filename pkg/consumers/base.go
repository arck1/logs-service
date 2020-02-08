package consumers


type ConsumerMessage struct {
	Meta    []byte
	Content []byte
}

type BaseConsumer interface {
	//insert consumed message into channel
	Listen(chan ConsumerMessage) error
	Stop()
}
