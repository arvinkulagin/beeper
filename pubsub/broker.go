package pubsub

import (
	"errors"
	"sync"
)

func deleteSubscriber(subscribers []chan []byte, subscriber chan []byte) []chan []byte {
	for i, v := range subscribers {
		if v == subscriber {
			return append(subscribers[:i], subscribers[i+1:]...)
		}
	}
	return subscribers
}

type Topic struct {
	subscribe   chan chan []byte
	unsubscribe chan chan []byte
	publish     chan []byte
	count       chan chan int
	kill        chan struct{}
}

func NewTopic() Topic {
	topic := Topic{
		subscribe:   make(chan chan []byte),
		unsubscribe: make(chan chan []byte),
		publish:     make(chan []byte),
		count:       make(chan chan int),
		kill:        make(chan struct{}),
	}
	go func() {
		subscribers := []chan []byte{}
		defer func() {
			for _, subscriber := range subscribers {
				close(subscriber)
			}
		}()
		for {
			select {
			case subscriber := <-topic.subscribe:
				subscribers = append(subscribers, subscriber)
			case subscriber := <-topic.unsubscribe:
				subscribers = deleteSubscriber(subscribers, subscriber)
				close(subscriber)
			case msg := <-topic.publish:
				for _, subscriber := range subscribers {
					subscriber <- msg
				}
			case c := <-topic.count:
				c <- len(subscribers)
			case <-topic.kill:
				return
			}
		}
	}()
	return topic
}

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

func (b *Broker) Subscribe(id string) (chan []byte, error) {
	subscriber := make(chan []byte) // Add buffer
	b.mutex.Lock()
	defer b.mutex.Unlock()
	topic, ok := b.topics[id]
	if !ok {
		return subscriber, errors.New("Topic does not exist")
	}
	topic.subscribe <- subscriber
	return subscriber, nil
}

func (b *Broker) Unsubscribe(id string, subscriber chan []byte) error {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	topic, ok := b.topics[id]
	if !ok {
		return errors.New("Topic does not exist")
	}
	topic.unsubscribe <- subscriber
	return nil
}

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