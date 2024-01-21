package eventbus

import (
	"chatmerger/internal/domain/model"
	"log"
	"sync"
)

type Event struct {
	Message          *model.Message
	DropSubscription *struct{}
}

type Subject = model.ID // client.ID

type Handler func(event Event) error

type EventBus struct {
	handlers map[Subject]Handler
	mu       *sync.Mutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[Subject]Handler),
		mu:       new(sync.Mutex),
	}
}

func (m *EventBus) Subscribe(subj Subject, handler Handler) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.handlers[subj] = handler
}

func (m *EventBus) Unsubscribe(subj Subject) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.handlers, subj)
}

func (m *EventBus) Publish(event Event, subjects ...Subject) error {
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
