package pkg

import (
	"log"
	"sync"
)

type Event interface {
	Name() string
}

type EventHandler interface {
	Handle(event Event)
}

type EventBus interface {
	Subscribe(eventName string, handler EventHandler)
	Publish(event Event)
}

type EventBusImpl struct {
	handlers map[string][]EventHandler
	rwMutex  sync.RWMutex
}

func (bus *EventBusImpl) Subscribe(eventName string, handler EventHandler) {
	bus.rwMutex.Lock()
	defer bus.rwMutex.Unlock()
	bus.handlers[eventName] = append(bus.handlers[eventName], handler)
	log.Printf("Subscribing to %s event", eventName)
}

func (bus *EventBusImpl) Publish(event Event) {
	bus.rwMutex.RLock()
	defer bus.rwMutex.RUnlock()
	for _, handler := range bus.handlers[event.Name()] {
		log.Printf("Publishing %s event", event.Name())
		go handler.Handle(event)
	}
}

func NewEventBus() EventBus {
	return &EventBusImpl{
		handlers: make(map[string][]EventHandler),
	}
}
