package server

import "github.com/greggolang/eventStoreClient"

type Topic struct {
	Name          string
	Listeners     *eventStoreClient.Store[Client]
	subscriberMap map[uint64]uint64
}

func newTopic(name string) *Topic {
	return &Topic{
		Name:          name,
		Listeners:     eventStoreClient.NewStore[Client](),
		subscriberMap: make(map[uint64]uint64),
	}
}

func (t *Topic) Subscribe(client *Client) {
	// Ignore the client subscribing to the topic it has already subscribed.
	if _, ok := t.subscriberMap[client.ConnID]; ok {
		return
	}

	id := t.Listeners.Put(client)
	t.subscriberMap[client.ConnID] = id
}

func (t *Topic) Unsubscribe(client *Client) {
	if subID, ok := t.subscriberMap[client.ConnID]; ok {
		t.Listeners.Delete(subID)
		delete(t.subscriberMap, client.ConnID)
	}
}

func (t *Topic) Publish(client *Client, data []byte) {
	t.Listeners.ForEach(func(_ uint64, client *Client) {
		client.Write(t.Name, data)
	})
}

type TopicStore struct {
	store    *eventStoreClient.Store[Topic]
	topicMap map[string]uint64
}

func NewTopicStore() *TopicStore {
	return &TopicStore{
		store:    eventStoreClient.NewStore[Topic](),
		topicMap: make(map[string]uint64),
	}
}

func (t *TopicStore) Subscribe(client *Client, name string) {
	var topic *Topic

	if id, ok := t.topicMap[name]; ok {
		topic, _ = t.store.Get(id)
	} else {
		topic = newTopic(name)
	}

	topic.Subscribe(client)
}

func (t *TopicStore) Unsubscribe(client *Client, name string) {
	if id, ok := t.topicMap[name]; ok {
		topic, _ := t.store.Get(id)
		topic.Unsubscribe(client)
	}
}

func (t *TopicStore) Publish(client *Client, name string, data []byte) {
	if id, ok := t.topicMap[name]; ok {
		topic, _ := t.store.Get(id)
		topic.Publish(client, data)
	}
}
