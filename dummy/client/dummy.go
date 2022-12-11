package dummy

import (
	"errors"
	"log"
	"net"

	eventeur "github.com/greggolang/eventStoreClient"
	proto "github.com/greggolang/eventStoreClient/dummy"
	"github.com/greggolang/eventStoreClient/dummy/tcp"
	"github.com/vmihailenco/msgpack"
)

type client struct {
	config Config
	msgC   chan tcp.Message
}

// NewDummyClient creates a client that demonstrates how to implement
// EventeurClient interface. *This should not be used in production*.
func NewDummyClient(config Config) eventeur.EventeurClient {
	return &client{
		config: config,
		msgC:   make(chan tcp.Message),
	}
}

func (c *client) Connect() error {
	for {
		conn, err := net.Dial("tcp", c.config.Url)
		if err != nil {
			if c.config.Reconnect {
				continue
			}
			log.Fatal(err)
		}
		transport := tcp.NewTransport(conn, c.handleMsg)

		go func() {
			for {
				msg, ok := <-c.msgC
				log.Print("got msg: ", string(msg))
				if !ok {
					break
				}

				err := transport.Write(msg)
				if err != nil {
					log.Print("unable to send msg: ", err)
				}
			}

			log.Print("client closed the send channel")
		}()

		transport.Spawn()
		break
	}

	return nil
}

func (c *client) handleMsg(w tcp.ResponseWriter, r tcp.Message) {
	// Convert raw bytes coming from connection to the wire message.
	mt, m, err := r.Unmarshal()
	if err != nil {
		log.Print(err)
		w.Close()
		return
	}

	// Message handling error. If our wire data was succesfully parsed,
	// but actions on that data failed, it will be logged at the end of
	// this function.
	var mErr error

	// Client only expects the to receive event messages from topics it's subscribed to.
	// And if a server decides to gracefully shutdown then it should send `Close` message.
	switch mt {
	case tcp.MsgOk:
	case tcp.MsgPublish:
		mErr = c.onEvent(m)
	case tcp.MsgUsubscribe:
		c.onClose(m)
	default:
		mErr = errors.New("unsupported message type")
		w.Close()
	}

	if mErr != nil {
		log.Print(mErr)
	}
}

func (c *client) onEvent(m []byte) error {
	return nil
}

func (c *client) onClose(m []byte) {}

// Publish implements eventeur.EventeurClient
func (c *client) Publish(topic string, data []byte) error {
	b, err := msgpack.Marshal(proto.PublishMsg{Topic: topic, Data: data})
	if err != nil {
		return err
	}

	c.msgC <- tcp.NewPublishMsg(b)
	return nil
}

// Subscribe implements eventeur.EventeurClient
func (c *client) Subscribe(topic string, callback eventeur.EventeurEventHandler) (eventeur.Subscription, error) {
	b, err := msgpack.Marshal(&proto.SubscribeMsg{Topic: topic})
	if err != nil {
		return nil, err
	}

	log.Print("subscribed")
	c.msgC <- tcp.NewSubscribeMsg(b)
	log.Print("moving")
	return newSubscription(c, topic), nil
}

// Unsubscribe implements eventeur.EventeurClient
func (c *client) Unsubscribe(topic string) error {
	b, err := msgpack.Marshal(&proto.UnsubscribeMsg{Topic: topic})
	if err != nil {
		return err
	}

	c.msgC <- tcp.NewPublishMsg(b)
	return nil
}

type subscription struct {
	client *client
	topic  string
}

func newSubscription(c *client, topic string) *subscription {
	return &subscription{client: c, topic: topic}
}

func (s *subscription) Unsubscribe() error {
	s.client.Unsubscribe(s.topic)
	return nil
}

// Close implements eventeur.EventeurClient
func (*client) Close() {}
