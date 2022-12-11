package server

type SubscribeMsg struct {
	Topic string
}

type UnsubscribeMsg struct {
	Topic string
}

type PublishMsg struct {
	Topic string
	Data []byte
}
