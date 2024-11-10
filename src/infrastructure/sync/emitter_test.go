package sync

import (
	"sync"
	"testing"
	"time"
)

func TestSubscribeAndSend(t *testing.T) {
	emitter := NewEmitter()

	ch1 := emitter.Subscribe()
	if ch1 == nil {
		t.Error("channel should not be nil")
	}

	ch2 := emitter.Subscribe()
	ch3 := emitter.Subscribe()

	// Send a message
	expectedMessage := "hello"
	go emitter.Send(expectedMessage)

	listeners := []<-chan string{ch1, ch2, ch3}
	wg := sync.WaitGroup{}

	for i := range listeners {
		go func(i int) {
			wg.Add(1)
			defer wg.Done()

			select {
			case msg := <-listeners[i]:
				if expectedMessage != msg {
					t.Error("message should be received correctly")
				}
			case <-time.After(1 * time.Second):
				t.Fatal("expected message not received")
			}
		}(i)
	}

	time.Sleep(1 * time.Second)

	wg.Wait()
}

func TestUnsubscribe(t *testing.T) {
	emitter := NewEmitter()

	ch := emitter.Subscribe()
	emitter.Unsubscribe(ch)

	// Try sending a message, but it shouldn't be received
	go emitter.Send("hello unsubscribe")

	time.Sleep(100 * time.Millisecond)

	_, open := <-ch
	if open {
		t.Error("expected the listener to not receive any message")
	}
}

func TestSubscribeTwice(t *testing.T) {
	emitter := NewEmitter()

	// First subscription should succeed
	ch1 := emitter.Subscribe()
	if ch1 == nil {
		t.Error("first subscription should be successful")
	}

	// Second subscription with the same ID should not fail
	ch2 := emitter.Subscribe()
	if ch2 == nil {
		t.Error("first subscription should be successful")
	}

	// Second channel must be different than the first
	if ch1 == ch2 {
		t.Error("channels must be different")
	}
}
