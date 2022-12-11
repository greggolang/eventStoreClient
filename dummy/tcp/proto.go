package tcp

import "errors"

// MessageType ...
type MessageType byte

// Message contains communication message type and it's body.
type Message []byte

const (
	// Newline representation in hex.
	// It's used to check for the end of a message that is comming through this transport.
	ByteLF = byte(0x0A)

	// MsgRequest is used by the client to request a challenge from the server.
	MsgConnect = MessageType(0x01)
	// MsgChallenge is sent by the server with a resource and required bit count for the proof.
	MsgSubscribe = MessageType(0x02)
	// MsgProof is sent by the client with a hashcash of a received resource.
	MsgPublish = MessageType(0x03)
	// MsgWords is sent by the server with the words of wisdom as a body.
	MsgUsubscribe = MessageType(0x04)
	// MsgError is a reserved message type for errors.
	MsgClose = MessageType(0x05)
	// MsgOk acknowledges that the party received and processed the message.
	MsgOk = MessageType(0x06)
)

// ErrorMessageTypeUnknown ...
var ErrorMessageTypeUnknown = errors.New("unknown message type")

// Validate checks if message type is known.
func (m MessageType) Validate() error {
	switch m {
	case MsgConnect, MsgSubscribe, MsgPublish, MsgUsubscribe, MsgClose, MsgOk:
		return nil
	default:
		return ErrorMessageTypeUnknown
	}
}

// Unmarshal parses the message and returns message type with a body if message type is valid.
func (m Message) Unmarshal() (MessageType, []byte, error) {
	mt := MessageType(m[0])

	// if a message is shorter than two bytes then it has no body.
	if len(m) < 2 {
		return mt, []byte{}, mt.Validate()
	}

	// removing the new line char when unmarshaling.
	return mt, m[1 : len(m)-1], mt.Validate()
}

// Marshal creates a new protocol message out of message type and body.
func (m *Message) Marshal(mt MessageType, body []byte) error {
	buf := append([]byte{byte(mt)}, body...)
	buf = append(buf, ByteLF)
	*m = Message(buf)
	return mt.Validate()
}

// NewMessage ...
func NewMessage(mt MessageType, body []byte) (Message, error) {
	m := Message{}
	err := m.Marshal(mt, body)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// NewConnectMsg ...
func NewConnectMsg() Message {
	m, _ := NewMessage(MsgConnect, []byte{})
	return m
}

// NewSubscribeMsg ...
func NewSubscribeMsg(body []byte) Message {
	m, _ := NewMessage(MsgSubscribe, body)
	return m
}

// NewPublishMsg ...
func NewPublishMsg(body []byte) Message {
	m, _ := NewMessage(MsgPublish, body)
	return m
}

// NewUnsubscribeMsg ...
func NewUnsubscribeMsg(body []byte) Message {
	m, _ := NewMessage(MsgUsubscribe, body)
	return m
}

// NewCloseMsg ...
func NewCloseMsg(body []byte) Message {
	m, _ := NewMessage(MsgClose, body)
	return m
}

// NewOkMsg ...
func NewOkMsg() Message {
	m, _ := NewMessage(MsgOk, []byte{})
	return m
}
