package events

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/suite"
	"sync"
	"testing"
	"time"
)

type TestEvent struct {
	Name     string
	Payload  interface{}
	Datetime time.Time
}

func (e *TestEvent) GetName() string {
	return e.Name
}

func (e *TestEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *TestEvent) SetPayload(p interface{}) {
	e.Payload = p
}

func (e *TestEvent) GetDatetime() time.Time {
	return e.Datetime
}

type TestEventHandler struct {
	ID string
}

func (eh *TestEventHandler) Handle(e EventInterface, group *sync.WaitGroup) {
}

func TestEventDispatcher_Register(t *testing.T) {
	eventDispatcher := NewEventDispatcher()
	event := &TestEvent{
		Name:     "Event 1",
		Payload:  "Test",
		Datetime: time.Now(),
	}
	eventHandler := &TestEventHandler{}
	eventHandler.Handle(event, nil)
	assert.Equal(t, 0, len(eventDispatcher.handlers[event.GetName()]))
	err := eventDispatcher.Register(event.GetName(), eventHandler)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(eventDispatcher.handlers[event.GetName()]))
	assert.Equal(t, eventHandler, eventDispatcher.handlers[event.GetName()][0])
}

func TestEventDispatcher_RegisterDuplicated(t *testing.T) {
	eventDispatcher := NewEventDispatcher()
	event := &TestEvent{
		Name:     "Event 1",
		Payload:  "Test",
		Datetime: time.Now(),
	}
	eventHandler := &TestEventHandler{}
	eventHandler.Handle(event, nil)
	err := eventDispatcher.Register(event.GetName(), eventHandler)
	assert.Nil(t, err)
	err = eventDispatcher.Register(event.GetName(), eventHandler)
	assert.Error(t, err)
}

func TestEventDispatcher_Clear(t *testing.T) {
	eventDispatcher := NewEventDispatcher()
	event := &TestEvent{
		Name:     "Event 1",
		Payload:  "Test",
		Datetime: time.Now(),
	}
	eventHandler := &TestEventHandler{}
	eventHandler.Handle(event, nil)
	eventDispatcher.Register(event.GetName(), eventHandler)
	eventDispatcher.Clear()
	assert.Equal(t, 0, len(eventDispatcher.handlers))
}

func TestEventDispatcher_Has(t *testing.T) {
	eventDispatcher := NewEventDispatcher()
	event := &TestEvent{
		Name:     "Event 1",
		Payload:  "Test",
		Datetime: time.Now(),
	}
	eventHandler := &TestEventHandler{
		ID: "1",
	}
	eventHandlerNotRegistered := &TestEventHandler{
		ID: "2",
	}
	eventHandler.Handle(event, nil)
	eventHandlerNotRegistered.Handle(event, nil)
	eventDispatcher.Register(event.GetName(), eventHandler)
	has := eventDispatcher.Has(event.GetName(), eventHandler)
	assert.True(t, has)
	has = eventDispatcher.Has(event.GetName(), eventHandlerNotRegistered)
	assert.False(t, has)
}

func TestEventDispatcher_Dispatch(t *testing.T) {
	eventDispatcher := NewEventDispatcher()
	event := &TestEvent{
		Name:     "Event 1",
		Payload:  "Test",
		Datetime: time.Now(),
	}
	eventHandler := &TestEventHandler{}
	eventHandler2 := &TestEventHandler{}
	eventHandler3 := &TestEventHandler{}
	eventHandler.Handle(event, nil)
	eventHandler2.Handle(event, nil)
	eventHandler3.Handle(event, nil)
	eventDispatcher.Register(event.GetName(), eventHandler)
	eventDispatcher.Register(event.GetName(), eventHandler2)
	eventDispatcher.Register(event.GetName(), eventHandler3)
	err := eventDispatcher.Dispatch(event)
	assert.Nil(t, err)
}

func TestEventDispatcher_Remove(t *testing.T) {
	eventDispatcher := NewEventDispatcher()
	event := &TestEvent{
		Name:     "Event 1",
		Payload:  "Test",
		Datetime: time.Now(),
	}
	eventHandler := &TestEventHandler{}
	eventHandler2 := &TestEventHandler{}
	eventHandler3 := &TestEventHandler{}
	eventHandler.Handle(event, nil)
	eventHandler2.Handle(event, nil)
	eventHandler3.Handle(event, nil)
	eventDispatcher.Register(event.GetName(), eventHandler)
	eventDispatcher.Register(event.GetName(), eventHandler2)
	eventDispatcher.Register(event.GetName(), eventHandler3)
	assert.True(t, eventDispatcher.Has(event.GetName(), eventHandler))
	assert.True(t, eventDispatcher.Has(event.GetName(), eventHandler2))
	assert.True(t, eventDispatcher.Has(event.GetName(), eventHandler3))
	err := eventDispatcher.Remove(event.GetName(), eventHandler2)
	assert.Nil(t, err)
	assert.False(t, eventDispatcher.Has(event.GetName(), eventHandler2))
	assert.True(t, eventDispatcher.Has(event.GetName(), eventHandler))
	assert.True(t, eventDispatcher.Has(event.GetName(), eventHandler3))
}
