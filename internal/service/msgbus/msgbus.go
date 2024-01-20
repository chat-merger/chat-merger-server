package msgbus

import (
	"chatmerger/internal/domain/model"
	"log"
	"sync"
)

type Event = model.Message // todo +onDropSubscriptionEvent
type Subject = model.ID    // client.ID

type Handler func(event Event) error

type MessagesBus struct {
	handlers map[Subject]Handler
	mu       sync.Mutex
}

func NewMessagesBus() *MessagesBus {
	return &MessagesBus{
		handlers: make(map[Subject]Handler),
		mu:       sync.Mutex{},
	}
}

func (m *MessagesBus) Subscribe(subj Subject, handler Handler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[subj] = handler
}

func (m *MessagesBus) Unsubscribe(subj Subject) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.handlers, subj)
}

func (m *MessagesBus) Publish(event Event, subjects ...Subject) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, s := range subjects {
		handler, ok := m.handlers[s]
		if !ok {
			log.Printf("[ERROR] handler with subj %s not found", s)
			continue
		}
		err := handler(event)
		if err != nil {
			log.Printf("[ERROR] %+v", err)
			// return faults.Wrap(err)
		}
	}
	return nil
}
