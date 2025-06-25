package main

import (
	"fmt"
	"sync"
)

type Message struct {
	ID     string `json:"-"` // uuid
	Topic  string `json:"-"`
	Points int    `json:"points"`
	UserID string `json:"-"`
}

type InMemoryMessageBroker struct {
	subscribers map[string][]chan Message
	lock        sync.RWMutex
}

func NewInMemoryMessageBroker() *InMemoryMessageBroker {
	return &InMemoryMessageBroker{
		subscribers: make(map[string][]chan Message),
	}
}

func (i *InMemoryMessageBroker) Subscribe(key string) chan Message {
	i.lock.RLock()
	defer i.lock.RUnlock()

	ch := make(chan Message, 10)
	i.subscribers[key] = append(i.subscribers[key], ch)
	return ch
}

func (i *InMemoryMessageBroker) Publish(key string, message Message) {
	i.lock.RLock()
	defer i.lock.RUnlock()

	subscribers := i.subscribers[key]
	for _, ch := range subscribers {
		select {
		case ch <- message:
			fmt.Printf("Message published for '%s': %v\n", key, message)
		default:
			fmt.Printf("Message NOT published for '%s': %v\n", key, message)
		}
	}
}
