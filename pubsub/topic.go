package pubsub

func deleteSubscriber(senders []Sender, sender Sender) []Sender {
	for i, v := range senders {
		if v.Enqueue() == sender.Enqueue() {
			return append(senders[:i], senders[i+1:]...)
		}
	}
	return senders
}

type Topic struct {
	subscribe   chan Sender
	unsubscribe chan Sender
	publish     chan []byte
	count       chan chan int
	kill        chan struct{}
}

func NewTopic() Topic {
	topic := Topic{
		subscribe:   make(chan Sender),
		unsubscribe: make(chan Sender),
		publish:     make(chan []byte),
		count:       make(chan chan int),
		kill:        make(chan struct{}),
	}
	go func() {
		senders := []Sender{}
		defer func() {
			for _, sender := range senders {
				sender.Close()
			}
		}()
		for {
			select {
			case sender := <-topic.subscribe:
				senders = append(senders, sender)
			case sender := <-topic.unsubscribe:
				senders = deleteSubscriber(senders, sender)
				sender.Close()
			case msg := <-topic.publish:
				for _, sender := range senders {
					sender.Send(msg)
				}
			case c := <-topic.count:
				c <- len(senders)
			case <-topic.kill:
				return
			}
		}
	}()
	return topic
}

func (t Topic) Subscribe() Receiver {
	sender, receiver := NewQueue(t.unsubscribe)
	t.subscribe <- sender
	return receiver
}