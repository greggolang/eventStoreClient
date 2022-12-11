package server

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/greggolang/eventStoreClient"
	proto "github.com/greggolang/eventStoreClient/dummy"
	"github.com/greggolang/eventStoreClient/dummy/tcp"
	"github.com/vmihailenco/msgpack"
)

type Server struct {
	config   Config
	clients  *eventeur.Store[Client]
	topics   *TopicStore
	messages *eventeur.Store[Message]
	queues   map[string]chan []byte
}

func NewDummyServer(config Config) *Server {
	return &Server{
		clients: eventeur.NewStore[Client](),
		topics:  NewTopicStore(),
		// TODO: implement events storage for replays.
		messages: eventeur.NewStore[Message](),
	}
}

// Publish implements eventeur.EventeurServer
func (s *Server) Publish(clientID uint64, topic string, data []byte) error {
	return errors.New("unimplemented")
}

// Subscribe implements eventeur.EventeurServer
func (s *Server) Subscribe(clientID uint64, topic string) error {
	return errors.New("unimplemented")
}

// Serve starts the TCP listener. In other words, this function
// runs in the background and accepts all the connections from
// remote clients.
//
// When a client connects, the server puts the information related
// to this client into the ClientStore and waits for incoming evets.
//
// The connection is also used for publishing events from
// server to that client (bidirectional connection).
func (s *Server) Serve(port int) error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go func() {
			tErr := s.handleTcp(conn)
			if tErr != nil {
				log.Print(err)
			}
		}()
	}
}

// handleTcp takes care of a low level communication that happens between
// server and client.
func (s *Server) handleTcp(conn net.Conn) error {
	var client Client
	id := s.clients.Put(&client)
	client.ConnID = id

	transport := tcp.NewTransport(conn, func(w tcp.ResponseWriter, r tcp.Message) {
		s.handleMsg(&client, w, r)
	})

	return transport.Spawn()
}

// handleMsg parses the low level messages that encapsulate event data.
// TODO: write response on every message.
func (s *Server) handleMsg(c *Client, w tcp.ResponseWriter, r tcp.Message) {
	// Convert raw bytes coming from connection to the wire message.
	mt, m, err := r.Unmarshal()
	if err != nil {
		log.Print("Error unmarshalling: ", err)
		w.Close()
		return
	}

	// Message handling error. If our wire data was succesfully parsed,
	// but actions on that data failed, it will be logged at the end of
	// this function.
	var mErr error

	log.Print("got msg: ", mt, string(m), err)
	switch mt {
	case tcp.MsgSubscribe:
		mErr = s.onSubscribe(c, m)
	case tcp.MsgUsubscribe:
		mErr = s.onUnsubscribe(c, m)
	case tcp.MsgPublish:
		mErr = s.onPublish(c, m)
	case tcp.MsgClose:
		s.onClose(c, m)
	default:
		mErr = errors.New("unsupported message type")
		w.Close()
	}

	if mErr != nil {
		log.Print("error when handling msg: ", mErr)
	}

	// Acknowledge the client that message was handled.
	w.Write(tcp.NewOkMsg())
	log.Print("msg processed")
}

func (s *Server) onSubscribe(c *Client, m []byte) error {
	var msg proto.SubscribeMsg
	err := msgpack.Unmarshal(m, &msg)
	if err != nil {
		return err
	}

	s.topics.Subscribe(c, msg.Topic)
	return nil
}

func (s *Server) onUnsubscribe(c *Client, m []byte) error {
	var msg proto.UnsubscribeMsg
	err := msgpack.Unmarshal(m, &msg)
	if err != nil {
		return err
	}

	s.topics.Unsubscribe(c, msg.Topic)
	return nil
}

func (s *Server) onPublish(c *Client, m []byte) error {
	var msg proto.PublishMsg
	err := msgpack.Unmarshal(m, &msg)
	if err != nil {
		return err
	}

	s.topics.Publish(c, msg.Topic, msg.Data)
	return nil
}

func (s *Server) onClose(c *Client, m []byte) {}

type Client struct {
	ConnID  uint64
	writeFn func([]byte) error
}

func (c *Client) Write(topic string, data []byte) error {
	b, err := msgpack.Marshal(&proto.PublishMsg{Topic: topic, Data: data})
	if err != nil {
		return err
	}

	return c.writeFn(b)
}

type Message struct{}
