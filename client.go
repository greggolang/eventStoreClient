package eventeur

type EventeurEventHandler func(clientID uint64, data []byte)

type Subscription interface {
	Unsubscribe() error
}

type EventeurClient interface {
	Connect() error
	Publish(topic string, data []byte) error
	Subscribe(topic string, callback EventeurEventHandler) (Subscription, error)
	Unsubscribe(topic string) error
	Close()
}
