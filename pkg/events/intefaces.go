package events

import (
	"sync"
	"time"
)

type EventInterface interface {
	GetName() string
	GetDatetime() time.Time
	GetPayload() interface{}
	SetPayload(p interface{})
}

type EventHandlerInterface interface {
	Handle(e EventInterface, group *sync.WaitGroup)
}

type EventDispatcherInterface interface {
	Register(eventName string, handler EventHandlerInterface) error
	Dispatch(e EventInterface) error
	Remove(eventName string, handler EventHandlerInterface) error
	Has(eventName string, handler EventHandlerInterface) bool
	Clear() error
}
