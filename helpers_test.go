package events

import (
	"fmt"
	"time"
)

// ExampleTicker shows how to use the Ticker helper to schedule events at regular intervals.
func ExampleTicker() {
	sim := NewSimulation()

	Ticker(sim, sim.Now, 3*time.Second, func(s *Simulation) {
		fmt.Println("Actor:", sim.Now)
	})

	whenToStop := sim.Now.Add(15 * time.Second)
	sim.RunUntil(whenToStop)

	// Output:
	// Actor: 0001-01-01 00:00:00 +0000 UTC
	// Actor: 0001-01-01 00:00:03 +0000 UTC
	// Actor: 0001-01-01 00:00:06 +0000 UTC
	// Actor: 0001-01-01 00:00:09 +0000 UTC
	// Actor: 0001-01-01 00:00:12 +0000 UTC
	// Actor: 0001-01-01 00:00:15 +0000 UTC
}
