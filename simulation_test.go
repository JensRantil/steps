package steps

import (
	"fmt"
	"testing"
	"time"
)

// ExampleSimulation shows how to use the Simulation type to schedule a single event and run a simulation until a given time.
func ExampleSimulation() {
	sim := NewSimulation()

	sim.Schedule(Event{When: sim.Now.Add(10 * time.Second), Action: func(s *Simulation) {
		fmt.Println("Actor 1:", sim.Now)
	}})
	sim.Schedule(Event{When: sim.Now.Add(time.Second), Action: func(s *Simulation) {
		fmt.Println("Actor 2:", sim.Now)
	}})

	sim.RunUntilDone()

	// Output:
	// Actor 2: 0001-01-01 00:00:01 +0000 UTC
	// Actor 1: 0001-01-01 00:00:10 +0000 UTC
}

func TestBasicSimulation(t *testing.T) {
	var timesCalled []time.Time

	sim := NewSimulation()

	whenToRun := sim.Now.Add(1 * time.Second)
	sim.Schedule(Event{When: whenToRun, Action: func(s *Simulation) {
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

	sim := NewSimulation()

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

func TestCancellingExistingEvent(t *testing.T) {
	sim := NewSimulation()

	id := sim.Schedule(Event{When: sim.Now.Add(time.Second), Action: func(s *Simulation) {
		fmt.Println("Actor 1:", sim.Now)
	}})

	if !sim.Cancel(id) {
		t.Errorf("expected event to be cancelled")
	}
}

func TestCancellingNonExistingEvent(t *testing.T) {
	sim := NewSimulation()

	id := sim.Schedule(Event{When: sim.Now.Add(time.Second), Action: func(s *Simulation) {
		fmt.Println("Actor 1:", sim.Now)
	}})

	nonExistingID := id + 1
	if sim.Cancel(nonExistingID) {
		t.Errorf("expected event to be cancelled")
	}
}
