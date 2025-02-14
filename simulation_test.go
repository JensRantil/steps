package events

import (
	"fmt"
	"testing"
	"time"
)

// ExampleSimulation shows how to use the Simulation type to schedule a single event and run a simulation until a given time.
func ExampleSimulation() {
	sim := &Simulation{}
	sim.Schedule(ScheduledEvent{When: sim.Now.Add(time.Second), Action: func(s *Simulation) {
		fmt.Println("Actor 1:", sim.Now)
	}})
	sim.RunUntil(sim.Now.Add(time.Second * 2))

	// Output:
	// Actor 1: 0001-01-01 00:00:01 +0000 UTC
}

func TestBasicSimulation(t *testing.T) {
	var timesCalled []time.Time

	sim := &Simulation{}

	whenToRun := sim.Now.Add(1 * time.Second)
	sim.Schedule(ScheduledEvent{When: whenToRun, Action: func(s *Simulation) {
		timesCalled = append(timesCalled, s.Now)
	}})

	whenToStop := sim.Now.Add(20 * time.Second)
	sim.RunUntil(whenToStop)

	if expectedTimesCalled := 1; len(timesCalled) != expectedTimesCalled {
		t.Errorf("expected %d times called, got %d", expectedTimesCalled, len(timesCalled))
	}

	if lastTime := timesCalled[len(timesCalled)-1]; !lastTime.Equal(whenToRun) {
		t.Errorf("last time called %s is not equal to %s", lastTime, whenToRun)
	}
}

func TestSimulationUntil(t *testing.T) {
	var timesCalled []time.Time

	sim := &Simulation{}

	Ticker(sim, sim.Now, 1*time.Second, func(s *Simulation) {
		timesCalled = append(timesCalled, s.Now)
	})
	Ticker(sim, sim.Now, 3*time.Second, func(s *Simulation) {
		timesCalled = append(timesCalled, s.Now)
	})
	Ticker(sim, sim.Now, 5*time.Second, func(s *Simulation) {
		timesCalled = append(timesCalled, s.Now)
	})

	whenToStop := sim.Now.Add(20 * time.Second)
	sim.RunUntil(whenToStop)

	if expectedTimesCalled := 33; len(timesCalled) != expectedTimesCalled {
		t.Errorf("expected %d times called, got %d", expectedTimesCalled, len(timesCalled))
	}

	if lastTime := timesCalled[len(timesCalled)-1]; lastTime.After(whenToStop) {
		t.Errorf("last time called %s is after %s", lastTime, whenToStop)
	}
}
