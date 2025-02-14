package events

import "time"

// Simulation runs a discrete event simulation.
type Simulation struct {
	// Now represents the current point in time in the simulation. It is not recommended to modify this value during a simulation.
	Now time.Time

	// nextID is incremented for each event scheduled to the simulation. It is used to sort events with the same time.
	nextID int

	// queue is the queue of future events to be processed.
	queue eventQueue
}

// Step advances the simulation by one time unit. It returns true if the simulation advanced, false if there were no events to process.
func (s *Simulation) Step() bool {
	if s.queue.Len() == 0 {
		return false
	}
	e := s.queue.Pop()
	if e.When.After(s.Now) {
		// Never allow s.Now to go backwards in time.
		s.Now = e.When
	}
	e.Action(s)
	return true
}

type Action func(*Simulation)

// Schedule adds an event to the simulation.
func (s *Simulation) Schedule(w time.Time, a Action) {
	s.queue.Push(scheduledEvent{Order: s.nextID, When: w, Action: a})
	s.nextID++
}

// RunUntil runs the simulation until the given time or there are no more events to process.
func (s *Simulation) RunUntil(until time.Time) {
	for {
		if s.queue.Len() == 0 {
			break
		}
		if s.queue.Peek().When.After(until) {
			// Don't process events after the given time.
			break
		}
		if !s.Step() {
			// Strictly speaking, this should never happen since we have the check for the queue length above. Better safe than sorry, though.
			break
		}
	}
}

// Ticker schedules an event to run at a regular interval. For more complex cronjobs etc., have a look at something like [1].
//
// [1]: https://pkg.go.dev/github.com/robfig/cron#Schedule
func Ticker(sim *Simulation, start time.Time, duration time.Duration, f func(s *Simulation)) {
	var nextRun func(s *Simulation)
	nextRun = func(s *Simulation) {
		f(s)
		s.Schedule(s.Now.Add(duration), nextRun)
	}

	// Schedule the first run.
	sim.Schedule(start, nextRun)
}
