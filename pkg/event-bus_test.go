package pkg

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"testing"
)

func TestEventBus(t *testing.T) {
}

type TestListener struct {
}

func (listener *TestListener) Handle(event Event) {
	log.Printf("Received event: %v", event)
}

const eventName = "TestEvent"

type TestEvent struct {
}

func (te *TestEvent) Name() string {
	return eventName
}

// test that Event Bus would be empty
func TestEventBusImpl(t *testing.T) {
	t.Parallel()

	testBus := NewEventBus()
	testBusImpl := testBus.(*EventBusImpl)

	if testBusImpl == nil {
		t.Fatal("NewEventBus() returned nil")
	}

	if len(testBusImpl.handlers) != 0 {
		t.Error("NewEventBus() did not return an empty list")
	}
}

// test that calling Subscribe adds the listener to the map
func TestEventBusImpl_Subscribe(t *testing.T) {
	t.Parallel()

	testBus := NewEventBus()
	testBusImpl := testBus.(*EventBusImpl)

	if testBusImpl == nil {
		t.Fatal("NewEventBus() returned nil")
	}

	// Random number between 1 and 5
	numEvents := rand.Intn(5) + 1

	for i := 0; i < numEvents; i++ {
		eventName := fmt.Sprintf("Event_%d", i)

		// Random number between 1 and 20
		numListeners := rand.Intn(20) + 1

		for j := 0; j < numListeners; j++ {
			testBusImpl.Subscribe(eventName, &TestListener{})
		}

		if len(testBusImpl.handlers[eventName]) != numListeners {
			t.Errorf(
				"Event %s: expected %v listeners, got %v",
				eventName,
				numListeners,
				len(testBusImpl.handlers[eventName]),
			)
		}
	}
}

type CountingListener struct {
	wg    *sync.WaitGroup
	mu    sync.Mutex
	count int
}

func (cl *CountingListener) Handle(e Event) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	defer cl.wg.Done()
	cl.count++
}

func (cl *CountingListener) Count() int {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	return cl.count
}

// test that when Published is called, all the listeners are called
func TestEventBusImpl_Publish(t *testing.T) {
	t.Parallel()
	testBus := NewEventBus()
	testBusImpl := testBus.(*EventBusImpl)

	if testBusImpl == nil {
		t.Fatal("NewEventBus() returned nil")
	}

	// Create multiple listeners
	numListeners := 5
	wg := sync.WaitGroup{}
	wg.Add(numListeners)
	listener := &CountingListener{wg: &wg}
	for i := 0; i < numListeners; i++ {
		testBusImpl.handlers[eventName] = append(testBusImpl.handlers[eventName], listener)
	}

	testEvent := &TestEvent{}
	testBusImpl.Publish(testEvent)
	wg.Wait()
	if numListeners != listener.Count() {
		t.Errorf("Expected %d Listeners to be called, %d called", numListeners, listener.Count())
	}
}
