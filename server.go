package eventeur

type EventeurServer interface {
	Publish(clientID uint64, topic string, data []byte) error
	Subscribe(clientID uint64, topic string) error
}
