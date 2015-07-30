package pubsub

import (
	"errors"
	"sync"
)

type Broker struct {
	topics map[string]Topic
	mutex  sync.Mutex
}

func NewBroker() *Broker {
	broker := Broker{
		topics: make(map[string]Topic),
		mutex:  sync.Mutex{},
	}
	return &broker
}

func (b *Broker) AddTopic(id string) error {
	topic := NewTopic()
	b.mutex.Lock()
	defer b.mutex.Unlock()
	_, ok := b.topics[id]
	if ok {
		return errors.New("Topic already exists")
	}
	b.topics[id] = topic
	return nil
}

func (b *Broker) DelTopic(id string) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	topic, ok := b.topics[id]
	if !ok {
		return errors.New("Topic does not exist")
	}
	topic.kill <- struct{}{}
	delete(b.topics, id)
	return nil
}

func (b *Broker) Topics() []string {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	ids := []string{}
	for id, _ := range b.topics {
		ids = append(ids, id)
	}
	return ids
}

func (b *Broker) Count(id string) (int, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	topic, ok := b.topics[id]
	if !ok {
		return 0, errors.New("Topic does not exist")
	}
	c := make(chan int)
	topic.count <- c
	return <-c, nil
}

func (b *Broker) Subscribe(id string) (Receiver, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	topic, ok := b.topics[id]
	if !ok {
		return Receiver{}, errors.New("Topic does not exist")
	}
	receiver := topic.Subscribe()
	return receiver, nil
}


// func (b *Broker) Unsubscribe(id string, subscriber chan []byte) error {
// 	b.mutex.Lock()
// 	defer b.mutex.Unlock()
// 	topic, ok := b.topics[id]
// 	if !ok {
// 		return errors.New("Topic does not exist")
// 	}
// 	topic.unsubscribe <- subscriber
// 	return nil
// }


func (b *Broker) Publish(id string, message []byte) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	topic, ok := b.topics[id]
	if !ok {
		return errors.New("Topic does not exist")
	}
	topic.publish <- message
	return nil
}