package pubsub

import (
	"testing"
	"sync"
)

func TestAddTopic(t *testing.T) {
	broker := NewBroker()
	err := broker.AddTopic("test")
	if err != nil {
		t.Error("Expected nil, got error:", err)
	}
	err = broker.AddTopic("test")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestDelTopic(t *testing.T) {
	broker := NewBroker()
	err := broker.AddTopic("test")
	if err != nil {
		t.Error("Expected nil, got error:", err)
	}
	err = broker.DelTopic("test")
	if err != nil {
		t.Error("Expected nil, got error:", err)
	}
	err = broker.DelTopic("test")
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestTopics(t *testing.T) {
	topics := []string{"test1", "test2", "test3", "test4", "test5"}
	broker := NewBroker()
	for _, topic := range topics {
		broker.AddTopic(topic)
	}
	if len(broker.Topics()) != len(topics) {
		t.Error("Wrong number of topics")
	}
}


func TestSubscribe(t *testing.T) {
	broker := NewBroker()
	_, err := broker.Subscribe("test")
	if err == nil {
		t.Error("Subscribe to nonexistent topic")
	}
	broker.AddTopic("test")
	_, err = broker.Subscribe("test")
	if err != nil {
		t.Error("Can not subscribe topic")
	}
	broker.DelTopic("test")
	_, err = broker.Subscribe("test")
	if err == nil {
		t.Error("Subscribe to deleted topic")
	}
}

func TestUnsubscribe(t *testing.T) {
	wg := sync.WaitGroup{}
	broker := NewBroker()
	broker.AddTopic("test")
	subscriber, _ := broker.Subscribe("test")
	wg.Add(1)
	go func(subscriber chan []byte) {
		_, ok := <-subscriber
		if !ok {
			t.Error("Channel closed")
		}
		wg.Done()
	}(subscriber)
	broker.Publish("test", []byte("Lorem ipsum dolor sit amet"))
	wg.Wait()
	broker.Unsubscribe("test", subscriber)
	wg.Add(1)
	go func(subscriber chan []byte) {
		_, ok := <-subscriber
		if ok {
			t.Error("Channel must be closed")
		}
		wg.Done()
	}(subscriber)
	broker.Publish("test", []byte("Lorem ipsum dolor sit amet"))
	wg.Wait()
}

func TestPublish(t *testing.T) {
	wg := sync.WaitGroup{}
	message := []byte("Lorem ipsum dolor sit amet")
	broker := NewBroker()
	err := broker.Publish("test", message)
	if err == nil {
		t.Error("Subscribe to nonexistent topic")
	}
	broker.AddTopic("test")
	subscriber, _ := broker.Subscribe("test")
	go func(subscriber chan []byte) {
		defer wg.Done()
		msg := <-subscriber
		if string(msg) != string(message) {
			t.Error("msg != " + string(message))
		}
	}(subscriber)
	wg.Add(1)
	err = broker.Publish("test", message)
	if err != nil {
		t.Error("Can not publish to topic")
	}
	wg.Wait()
	broker.DelTopic("test")
	err = broker.Publish("test", message)
	if err == nil {
		t.Error("Publish to nonexistent topic")
	}
}

func TestCount(t *testing.T) {
	numberOfSubscribers := 10
	broker := NewBroker()
	_, err := broker.Count("test")
	if err == nil {
		t.Error("Count in nonexistent topic")
	}
	broker.AddTopic("test")
	for i := 0; i < 10; i++ {
		broker.Subscribe("test")
	}
	count, _ := broker.Count("test")
	if numberOfSubscribers != count {
		t.Error("Wrong number of subscribers")
	}
}