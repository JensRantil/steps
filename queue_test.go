package events

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
		{ID: 1, Event: ScheduledEvent{When: now.Add(time.Second), Action: nil}},
		{ID: 2, Event: ScheduledEvent{When: now.Add(time.Second * 2), Action: nil}},
		{ID: 3, Event: ScheduledEvent{When: now.Add(time.Second * 2), Action: nil}},
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
