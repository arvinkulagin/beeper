package pubsub

import (
	"errors"
)

func NewQueue(unsubscribe chan Sender) (Sender, Receiver) {
	sender := Sender{
		enqueue: make(chan []byte),
	}
	receiver := Receiver{
		dequeue: make(chan []byte),
		kill: make(chan struct{}),
	}
	go func() {
		defer close(receiver.dequeue)
		messages := [][]byte{}
		for {
			if len(messages) == 0 {
				select {
				case message, ok := <-sender.enqueue:
					if !ok {
						return
					}
					messages = append(messages, message)
				case <-receiver.kill:
					unsubscribe <- sender
					return
				}
			} else {
				select {
				case message := <-sender.enqueue:
					messages = append(messages, message)
				case receiver.dequeue <- messages[0]:
					messages[0] = nil
					messages = messages[1:]
				case <-receiver.kill:
					unsubscribe <- sender
					return
				}
			}
		}
	}()
	return sender, receiver
}

type Sender struct {
	enqueue chan []byte
}

func (s Sender) Enqueue() chan []byte {
	return s.enqueue
}

func (s Sender) Send(message []byte) {
	s.enqueue <- message
}

func (s Sender) Close() {
	close(s.enqueue)
}

type Receiver struct {
	dequeue chan []byte
	kill chan struct{}
}

func (r Receiver) Receive() ([]byte, error) {
	message, ok := <-r.dequeue
	if !ok {
		return message, errors.New("Queue was closed")
	}
	return message, nil
}

func (r Receiver) Close() {
	r.kill <- struct{}{}
}