package steps

import (
	"container/heap"
	"fmt"
	"time"
)

// Event represents an event that will be processed as soon as possible after a specific time.
type Event struct {
	// When is the time at which the event should be processed as measured from the start of the simulation.
	When time.Time

	// Action is the function to call when the event is to be processed.
	Action Action
}

// scheduledEvent represents a future scheduledEvent in the simulation.
type scheduledEvent struct {
	ID    EventID
	Event Event
}

// String returns a string representation of the event.
func (e scheduledEvent) String() string {
	return fmt.Sprintf("scheduledEvent{ID: %d, When: %s}", e.ID, e.Event.When)
}

// eventQueue is a type-safe heap of events. Events with the same time are sorted by order. Otherwise, they are sorted by time, smallest first.
type eventQueue struct {
	heap eventsHeap
}

// newEventQueue creates a new event queue.
func newEventQueue() *eventQueue {
	return &eventQueue{
		heap: eventsHeap{
			IndexByID: make(map[EventID]int),
		},
	}
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
	return q.heap.Events[0]
}

// Remove removes an event from the queue. Returns true if the event was found and removed, false otherwise.
func (q *eventQueue) Remove(id EventID) bool {
	index, found := q.heap.IndexByID[id]
	if !found {
		return false
	}
	heap.Remove(&q.heap, index)
	return true
}

// Len returns the number of events in the queue.
func (q *eventQueue) Len() int {
	return q.heap.Len()
}

// A heap of events. Events with the same time are sorted by order. Otherwise, they are sorted by time, smallest first.
type eventsHeap struct {
	Events    []scheduledEvent
	IndexByID map[EventID]int
}

func (h eventsHeap) Len() int {
	sliceLen := len(h.Events)
	if sliceLen != len(h.IndexByID) {
		panic(fmt.Sprintf("len(h.Events) != len(h.IndexByID): %d != %d", sliceLen, len(h.IndexByID)))
	}
	return sliceLen
}
func (h eventsHeap) Less(i, j int) bool {
	if h.Events[i].Event.When == h.Events[j].Event.When {
		return h.Events[i].ID < h.Events[j].ID
	}
	return h.Events[i].Event.When.Before(h.Events[j].Event.When)
}
func (h eventsHeap) Swap(i, j int) {
	h.Events[i], h.Events[j] = h.Events[j], h.Events[i]

	h.IndexByID[h.Events[i].ID] = i
	h.IndexByID[h.Events[j].ID] = j
}

func (h *eventsHeap) Push(xUntyped any) {
	x := xUntyped.(scheduledEvent)
	if _, found := h.IndexByID[x.ID]; found {
		panic(fmt.Sprintf("event with ID %d already exists", x.ID))
	}
	h.Events = append(h.Events, x)
	h.IndexByID[x.ID] = len(h.Events) - 1
}

func (h *eventsHeap) Pop() any {
	n := len(h.Events)
	x := h.Events[n-1]
	h.Events = h.Events[0 : n-1]
	delete(h.IndexByID, x.ID)
	return x
}
