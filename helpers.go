package events

import "time"

// Ticker schedules an event to run at a regular interval. For more complex cronjobs etc., have a look at something like [1].
//
// [1]: https://pkg.go.dev/github.com/robfig/cron#Schedule
func Ticker(sim *Simulation, start time.Time, duration time.Duration, f func(s *Simulation)) {
	var nextRun func(s *Simulation)
	nextRun = func(s *Simulation) {
		f(s)
		s.Schedule(ScheduledEvent{When: s.Now.Add(duration), Action: nextRun})
	}

	// Schedule the first run.
	sim.Schedule(ScheduledEvent{When: start, Action: nextRun})
}
