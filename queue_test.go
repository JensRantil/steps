package steps

import (
	"math/rand/v2"
	"slices"
	"testing"
	"time"
)

func TestQueueTimeBasedOrdering(t *testing.T) {
	var now time.Time

	// These events are in order of schedule priority.
	events := []scheduledEvent{
		{ID: 1, Event: Event{When: now.Add(time.Second), Action: nil}},
		{ID: 2, Event: Event{When: now.Add(time.Second * 2), Action: nil}},
		{ID: 3, Event: Event{When: now.Add(time.Second * 2), Action: nil}},
	}

	// Subtest that adds the events to the queue and then pops them in order.
	t.Run("Add in order", func(t *testing.T) {
		queue := newEventQueue()
		for _, e := range events {
			queue.Push(e)
		}

		for _, e := range events {
			if popped := queue.Pop(); !eventsAreEqual(popped, e) {
				t.Errorf("Expected %v, got %v", e, popped)
			}
		}
	})
	t.Run("Add in reverse order", func(t *testing.T) {
		reversedEvents := slices.Clone(events)
		slices.Reverse(reversedEvents)

		queue := newEventQueue()
		for _, e := range reversedEvents {
			queue.Push(e)
		}

		for _, e := range events {
			if popped := queue.Pop(); !eventsAreEqual(popped, e) {
				t.Errorf("Expected %v, got %v", e, popped)
			}
		}
	})

	t.Run("Add in random order", func(t *testing.T) {
		randomEvents := slices.Clone(events)

		testIterations := 100
		for range testIterations {
			rand.Shuffle(len(randomEvents), func(i, j int) {
				randomEvents[i], randomEvents[j] = randomEvents[j], randomEvents[i]
			})

			queue := newEventQueue()
			for _, e := range randomEvents {
				queue.Push(e)
			}

			for _, e := range events {
				if popped := queue.Pop(); !eventsAreEqual(popped, e) {
					t.Errorf("Expected %v, got %v", e, popped)
				}
			}
		}
	})
}

func eventsAreEqual(a, b scheduledEvent) bool {
	return a.ID == b.ID && a.Event.When == b.Event.When
}

func TestQueueLen(t *testing.T) {
	queue := newEventQueue()
	queue.Push(scheduledEvent{ID: 1, Event: Event{When: time.Now().Add(time.Second), Action: nil}})
	queue.Push(scheduledEvent{ID: 2, Event: Event{When: time.Now().Add(time.Second * 2), Action: nil}})
	queue.Push(scheduledEvent{ID: 3, Event: Event{When: time.Now().Add(time.Second * 2), Action: nil}})

	if queue.Len() != 3 {
		t.Errorf("expected queue length 3, got %d", queue.Len())
	}
}

func TestQueueRemoveOfExistingEvent(t *testing.T) {
	queue := newEventQueue()

	event := scheduledEvent{ID: 42, Event: Event{When: time.Now().Add(time.Second), Action: nil}}
	queue.Push(event)

	if removed := queue.Remove(event.ID); !removed {
		t.Errorf("expected event to be removed")
	}

	if queue.Len() != 0 {
		t.Errorf("expected queue length 2, got %d", queue.Len())
	}
}

func TestQueueRemoveOfMissingEvent(t *testing.T) {
	queue := newEventQueue()

	event := scheduledEvent{ID: 42, Event: Event{When: time.Now().Add(time.Second), Action: nil}}
	queue.Push(event)

	missingID := event.ID - 1
	if removed := queue.Remove(missingID); removed {
		t.Errorf("expected event to not be removed")
	}

	if queue.Len() != 1 {
		t.Errorf("expected queue length 2, got %d", queue.Len())
	}
}
