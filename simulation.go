package events

import "time"

type ScheduledEventID int

// Simulation runs a discrete event simulation.
type Simulation struct {
	// Now represents the current point in time in the simulation. It is not recommended to modify this value during a simulation.
	Now time.Time

	// nextID is incremented for each event scheduled to the simulation. It is used to sort events with the same time.
	nextID ScheduledEventID

	// queue is the queue of future events to be processed.
	queue *eventQueue
}

// NewSimulation creates a new simulation.
func NewSimulation() *Simulation {
	// Currently, the zero value of Simulation is a valid simulation. However, this function exists to this library a bit more forward compatible in case zero values are no longer valid.
	return &Simulation{queue: newEventQueue()}
}

// Step advances the simulation by one time unit. It returns true if the simulation advanced, false if there were no events to process.
func (s *Simulation) Step() bool {
	if s.queue.Len() == 0 {
		return false
	}
	e := s.queue.Pop()
	if e.Event.When.After(s.Now) {
		// Never allow s.Now to go backwards in time.
		s.Now = e.Event.When
	}
	e.Event.Action(s)
	return true
}

type Action func(*Simulation)

// Schedule adds an event to the simulation.
func (s *Simulation) Schedule(e ScheduledEvent) ScheduledEventID {
	id := s.nextID
	s.queue.Push(scheduledEvent{ID: id, Event: e})
	s.nextID++
	return id
}

// Cancel cancels an event scheduled to the simulation. Returns true if the event was found and cancelled, false if the event was not found (never scheduled, or it was already executed).
func (s *Simulation) Cancel(id ScheduledEventID) bool {
	return s.queue.Remove(id)
}

// RunUntil runs the simulation until the given time or there are no more events to process.
func (s *Simulation) RunUntil(until time.Time) {
	for {
		if s.queue.Len() == 0 {
			break
		}
		if s.queue.Peek().Event.When.After(until) {
			// Don't process events after the given time.
			break
		}
		if !s.Step() {
			// Strictly speaking, this should never happen since we have the check for the queue length above. Better safe than sorry, though.
			break
		}
	}
}

// RunUntilDone runs the simulation until there are no more events to process.
func (s *Simulation) RunUntilDone() {
	for s.Step() {
		// Deliberately left empty.
	}
}
