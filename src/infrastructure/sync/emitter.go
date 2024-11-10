package sync

import (
	"sync"
)

const channelSize = 500

// Emitter allows communication between multiple producers and listeners.
type Emitter struct {
	mu        sync.RWMutex
	listeners map[chan string]bool
}

func NewEmitter() *Emitter {
	var s = &Emitter{
		listeners: make(map[chan string]bool),
	}

	return s
}

// Send allows goroutines to send messages over the emitter. Non-blocking sends are ensured.
func (s *Emitter) Send(msg string) {
	// we must protect the listeners map from changes during message transmission
	s.mu.RLock()
	defer s.mu.RUnlock()

	for listener := range s.listeners {
		if len(listener) <= channelSize {
			listener <- msg
		}
		// else message is discarded
	}
}

func (e *Emitter) Subscribe() chan string {
	e.mu.Lock()
	defer e.mu.Unlock()
	ch := make(chan string, channelSize)
	e.listeners[ch] = true

	return ch
}

func (e *Emitter) Unsubscribe(listener chan string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, ok := e.listeners[listener]; ok {
		close(listener)
	}

	delete(e.listeners, listener)
}
