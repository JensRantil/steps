package events

import (
	"testing"
	"time"
)

func TestBasicSimulation(t *testing.T) {
	var timesCalled []time.Time

	sim := &Simulation{}

	whenToRun := sim.Now.Add(1 * time.Second)
	sim.Schedule(whenToRun, func(s *Simulation) {
		timesCalled = append(timesCalled, s.Now)
	})

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
