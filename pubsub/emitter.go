package pubsub

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// EventData represents the data to pass, it should have a copy function
type EventData interface {
	Copy() interface{}
}

// EventHandler handles event
type EventHandler func(data EventData)

// DeregisterFunc is a function to deregister from event
type DeregisterFunc func()

// EventSubscriber is able to subscribe on events
type EventSubscriber interface {
	On(event string, handler EventHandler) DeregisterFunc
	Once(event string, handler EventHandler)
	Channel(event string) (<-chan EventData, DeregisterFunc)
}

// EventPublisher is able to notify events
type EventPublisher interface {
	Notify(event string, data EventData)
	Clear(event string)
}

// Emitter is managing events
type Emitter interface {
	EventSubscriber
	EventPublisher
}

// eventHandlers represents a map of event handlers, attached to some event
type eventHandlers map[string]EventHandler

// emitter implements Emitter
type emitter struct {
	handlers map[string]eventHandlers
	mut      sync.RWMutex
}

// NewEmitter creates a new instance of emitter
func NewEmitter() Emitter {
	return &emitter{
		handlers: map[string]eventHandlers{},
		mut:      sync.RWMutex{},
	}
}

// Channel abstracts Emitter.On() to event as a channel, by listening on event and
func (e *emitter) Channel(event string) (<-chan EventData, DeregisterFunc) {
	running := uint32(1)
	cn := make(chan EventData)
	ctx, cancel := context.WithCancel(context.Background())
	deregister := e.On(event, func(data EventData) {
		select {
		case <-ctx.Done():
			return
		default:
			if ctx.Err() != nil {
				return
			}
			if atomic.LoadUint32(&running) > uint32(0) {
				cn <- data
			}
		}
	})
	return cn, func() {
		atomic.StoreUint32(&running, uint32(0))
		cancel()
		deregister()
		go close(cn)
	}
}

// Once will call handler only once
func (e *emitter) Once(event string, handler EventHandler) {
	var once sync.Once
	var deregister DeregisterFunc
	deregister = e.On(event, func(data EventData) {
		once.Do(func() {
			go deregister()
			handler(data)
		})
	})
}

// On register to event, returns a function for de-registration
func (e *emitter) On(event string, handler EventHandler) DeregisterFunc {
	e.mut.Lock()
	defer e.mut.Unlock()

	var handlers eventHandlers
	var ok bool
	if handlers, ok = e.handlers[event]; !ok {
		handlers = make(eventHandlers)
	}
	h := sha256.Sum256([]byte(fmt.Sprintf("%s:%d", time.Now().String(), len(handlers))))
	hid := hex.EncodeToString(h[:])
	handlers[hid] = handler
	e.handlers[event] = handlers

	return func() {
		e.mut.Lock()
		defer e.mut.Unlock()

		if _handlers, ok := e.handlers[event]; ok {
			delete(_handlers, hid)
			e.handlers[event] = _handlers
		}
	}
}

// Notify notifies on event
func (e *emitter) Notify(event string, data EventData) {
	e.mut.RLock()
	defer e.mut.RUnlock()

	handlers, ok := e.handlers[event]
	if !ok {
		return
	}
	for _, handler := range handlers {
		edata := data.Copy().(EventData)
		go handler(edata)
	}
}

// Clear handlers for event
func (e *emitter) Clear(event string) {
	e.mut.Lock()
	defer e.mut.Unlock()

	delete(e.handlers, event)
}

// countHandlers returns the number of handlers active on the given event
func (e *emitter) countHandlers(event string) int {
	e.mut.Lock()
	defer e.mut.Unlock()

	if handlers, ok := e.handlers[event]; ok {
		return len(handlers)
	}
	return 0
}