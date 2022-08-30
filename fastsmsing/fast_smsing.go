package fastsmsing

import "sync"

func NewInMemoryClient() FastSmsingClient {
	return &inMemoryClient{subscribers: make([]chan map[MessageID]MessageStatus, 0), lock: sync.RWMutex{}}
}

type inMemoryClient struct {
	lock        sync.RWMutex
	subscribers []chan map[MessageID]MessageStatus
}

func (c *inMemoryClient) Send(messages []Message) error {
	c.confirmMessages(messages)
	c.markAsDelivered(messages)
	return nil
}

func (c *inMemoryClient) confirmMessages(messages []Message) {
	toConfirm := map[MessageID]MessageStatus{}
	for _, message := range messages {
		if len(message.MessageID) > 0 {
			toConfirm[message.MessageID] = CONFIRMED
		}
	}
	for _, subscriber := range c.subscribers {
		subscriber <- toConfirm
	}
}

func (c *inMemoryClient) markAsDelivered(messages []Message) {
	toDeliver := map[MessageID]MessageStatus{}
	for _, message := range messages {
		if len(message.MessageID) > 0 {
			toDeliver[message.MessageID] = DELIVERED
		}
	}
	for _, subscriber := range c.subscribers {
		subscriber <- toDeliver
	}
}

func (c *inMemoryClient) Subscribe(subscriber chan map[MessageID]MessageStatus) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.subscribers = append(c.subscribers, subscriber)
}

func (c *inMemoryClient) Stop() {
	c.lock.Lock()
	defer c.lock.Unlock()
	
	for _, subscriber := range c.subscribers {
		close(subscriber)
	}
}
