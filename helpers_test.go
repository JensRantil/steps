package events

import (
	"fmt"
	"time"
)

// ExampleTicker shows how to use the Ticker helper to schedule events at regular intervals.
func ExampleTicker() {
	sim := NewSimulation()

	sim.Schedule(ScheduledEvent{When: sim.Now.Add(time.Second), Action: func(s *Simulation) {
		fmt.Println("Actor 1:", sim.Now)
	}})
	Ticker(sim, sim.Now, 1*time.Second, func(s *Simulation) {
		fmt.Println("Actor 1:", sim.Now)
	})
	Ticker(sim, sim.Now, 3*time.Second, func(s *Simulation) {
		fmt.Println("Actor 2:", sim.Now)
	})
	Ticker(sim, sim.Now, 5*time.Second, func(s *Simulation) {
		fmt.Println("Actor 3:", sim.Now)
	})

	whenToStop := sim.Now.Add(5 * time.Second)
	sim.RunUntil(whenToStop)

	// Output:
	// Actor 1: 0001-01-01 00:00:00 +0000 UTC
	// Actor 2: 0001-01-01 00:00:00 +0000 UTC
	// Actor 3: 0001-01-01 00:00:00 +0000 UTC
	// Actor 1: 0001-01-01 00:00:01 +0000 UTC
	// Actor 1: 0001-01-01 00:00:01 +0000 UTC
	// Actor 1: 0001-01-01 00:00:02 +0000 UTC
	// Actor 2: 0001-01-01 00:00:03 +0000 UTC
	// Actor 1: 0001-01-01 00:00:03 +0000 UTC
	// Actor 1: 0001-01-01 00:00:04 +0000 UTC
	// Actor 3: 0001-01-01 00:00:05 +0000 UTC
	// Actor 1: 0001-01-01 00:00:05 +0000 UTC
}
