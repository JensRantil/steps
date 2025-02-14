package events

import (
	"container/heap"
	"fmt"
	"time"
)

// scheduledEvent represents a future scheduledEvent in the simulation.
type scheduledEvent struct {
	Order int

	// When is the time at which the event should be processed as measured from the start of the simulation.
	When time.Time

	// Action is the function to call when the event is to be processed.
	Action func(*Simulation)
}

// String returns a string representation of the event.
func (e scheduledEvent) String() string {
	return fmt.Sprintf("Event{Order: %d, When: %s}", e.Order, e.When)
}

// eventQueue is a type-safe heap of events. Events with the same time are sorted by order. Otherwise, they are sorted by time, smallest first.
type eventQueue struct {
	heap eventsHeap
}

// newEventQueue creates a new event queue.
func newEventQueue() *eventQueue {
	return &eventQueue{}
}

// Push adds an event to the queue.
func (q *eventQueue) Push(e scheduledEvent) {
	heap.Push(&q.heap, e)
}

// Pop removes the next event from the queue.
func (q *eventQueue) Pop() scheduledEvent {
	return heap.Pop(&q.heap).(scheduledEvent)
}

// Peek returns the next event in the queue without removing it.
func (q *eventQueue) Peek() scheduledEvent {
	return q.heap[0]
}

// Len returns the number of events in the queue.
func (q *eventQueue) Len() int {
	return q.heap.Len()
}

// A heap of events. Events with the same time are sorted by order. Otherwise, they are sorted by time, smallest first.
type eventsHeap []scheduledEvent

func (h eventsHeap) Len() int { return len(h) }
func (h eventsHeap) Less(i, j int) bool {
	if h[i].When == h[j].When {
		return h[i].Order < h[j].Order
	}
	return h[i].When.Before(h[j].When)
}
func (h eventsHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *eventsHeap) Push(x any) {
	*h = append(*h, x.(scheduledEvent))
}

func (h *eventsHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
