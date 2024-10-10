package events

import (
	"errors"
	"sync"
)

type EventDispatcher struct {
	handlers map[string][]EventHandlerInterface
}

func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandlerInterface),
	}
}

func (ed *EventDispatcher) Register(eventName string, handler EventHandlerInterface) error {
	if _, alreadyIn := ed.handlers[eventName]; alreadyIn {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return errors.New("eventHandler already registered")
			}
		}
	}
	ed.handlers[eventName] = append(ed.handlers[eventName], handler)
	return nil
}

func (ed *EventDispatcher) Clear() error {
	ed.handlers = make(map[string][]EventHandlerInterface)
	return nil
}

func (ed *EventDispatcher) Has(eventName string, handler EventHandlerInterface) bool {
	if _, exists := ed.handlers[eventName]; exists {
		for _, h := range ed.handlers[eventName] {
			if h == handler {
				return true
			}
		}
	}
	return false
}

func (ed *EventDispatcher) Dispatch(e EventInterface) error {
	handler, exists := ed.handlers[e.GetName()]
	if !exists {
		return errors.New("event not found")
	}
	group := &sync.WaitGroup{}
	for _, h := range handler {
		group.Add(1)
		go h.Handle(e, group)
	}
	group.Wait()
	return nil
}

func (ed *EventDispatcher) Remove(eventName string, handler EventHandlerInterface) error {
	if _, exists := ed.handlers[eventName]; !exists {
		return errors.New("event not found")
	}
	for idx, h := range ed.handlers[eventName] {
		if h == handler {
			ed.handlers[eventName] = append(ed.handlers[eventName][:idx], ed.handlers[eventName][idx+1:]...)
			return nil
		}
	}
	return nil
}
