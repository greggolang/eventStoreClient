package tcp

import (
	"bufio"
	"errors"
	"net"
	"sync"
)

var (
	ErrorSendChannelClosed = errors.New("send channel is closed")
	ErrorWriteToClosed     = errors.New("writing to a closed transport")
)

type ResponseWriter interface {
	Write(Message) error
	Close()
}

type HandleFunc func(ResponseWriter, Message)

type Transport struct {
	conn      net.Conn
	sendC     chan []byte
	stopC     chan struct{}
	handler   HandleFunc
	closeOnce sync.Once
	closed    bool
	mu        sync.RWMutex
}

func NewTransport(conn net.Conn, handler HandleFunc) *Transport {
	return &Transport{
		conn:    conn,
		handler: handler,
		sendC:   make(chan []byte),
		stopC:   make(chan struct{}),
	}
}

func (t *Transport) Spawn() error {
	defer t.Close()
	errC := make(chan error)

	go func() {
		t.read(errC)
	}()

	err := <-errC
	return err
}

func (t *Transport) Close() {
	t.closeOnce.Do(func() {
		t.mu.Lock()
		defer t.mu.Unlock()

		t.closed = true
		close(t.stopC)
		t.conn.Close()
	})
}

func (t *Transport) Write(msg Message) error {
	if t.IsClosed() {
		return ErrorWriteToClosed
	}

	_, err := t.conn.Write(msg)
	return err
}

func (t *Transport) IsClosed() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.closed
}

func (t *Transport) read(errC chan<- error) {
	for {
		select {
		case <-t.stopC:
			return
		default:
			bytes, err := bufio.NewReader(t.conn).ReadBytes(ByteLF)
			if err != nil {
				errC <- err
				return
			}

			t.handler(t, Message(bytes))
		}
	}
}
