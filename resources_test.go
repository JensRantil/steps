package steps

import (
	"fmt"
	"testing"
	"time"
)

func TestConditionSignal(t *testing.T) {
	// Given
	s := NewSimulation()
	c := NewCondition(s)
	a1 := &testAction{}
	a2 := &testAction{}
	c.Wait(a1.Execute)
	c.Wait(a2.Execute)
	s.RunUntilDone()
	if a1.executed {
		t.Error("a1 was executed before triggered")
	}
	if a2.executed {
		t.Error("a2 was executed before triggered")
	}

	// When
	c.Signal()
	s.RunUntilDone()

	// Then
	if !a1.executed {
		t.Error("a1 was not executed")
	}
	if a2.executed {
		t.Error("a2 was executed, but only signaled one action to be triggered")
	}
}

func TestConditionBroadcast(t *testing.T) {
	// Given
	s := NewSimulation()
	c := NewCondition(s)
	a1 := &testAction{}
	a2 := &testAction{}
	c.Wait(a1.Execute)
	c.Wait(a2.Execute)
	s.RunUntilDone()
	if a1.executed {
		t.Error("a1 was executed before triggered")
	}
	if a2.executed {
		t.Error("a2 was executed before triggered")
	}

	// When
	c.Broadcast()
	s.RunUntilDone()

	// Then
	if !a1.executed {
		t.Error("a1 was not executed but we broadcasted all actions to be triggered")
	}
	if !a2.executed {
		t.Error("a2 was not executed but we broadcasted all actions to be triggered")
	}
}

type testAction struct {
	executed bool
}

func (a *testAction) Execute(*Simulation) {
	a.executed = true
}

func TestBinarySemaphore(t *testing.T) {
	testCountingSemaphore(t, 1)
}

func TestCountingSemaphore(t *testing.T) {
	testCountingSemaphore(t, 4)
}

func testCountingSemaphore(t *testing.T, max int) {
	sim := NewSimulation()
	sem := NewCountingSemaphore(sim, max)

	maxParallelism := 0
	running := 0
	executions := 0

	nbrOfExecutions := 10 * max
	for range nbrOfExecutions {
		sem.Acquire(func(sim *Simulation) {
			running++
			if running > maxParallelism {
				maxParallelism = running
			}

			executions++

			sim.Schedule(Event{When: sim.Now.Add(10 * time.Second), Action: func(sim *Simulation) {
				sem.Release()
				running--
			}})
		})
	}
	sim.RunUntilDone()

	if maxParallelism != max {
		t.Errorf("maxParallel was %d, but should have been %d", maxParallelism, max)
	}

	if executions != nbrOfExecutions {
		t.Errorf("executions was %d, but should have been %d", executions, nbrOfExecutions)
	}
}

// ExampleCondition demonstrates how to use the Condition to synchronize actions. It simulates processing 100 items, rate-limited to one item per second.
func ExampleCondition() {
	sim := NewSimulation()
	c := NewCondition(sim)

	itemsToProcess := 10
	for range itemsToProcess {
		c.Wait(func(sim *Simulation) {
			fmt.Println(sim.Now, "Processing...")
		})
	}
	Ticker(sim, sim.Now, time.Second, func(sim *Simulation) {
		c.Signal()
	})

	// Deliberately not using sim.RunUntilDone() here since the Ticker will run indefinitely.
	sim.RunUntil(sim.Now.Add(time.Duration(2*itemsToProcess) * time.Second))

	// Output:
	// 0001-01-01 00:00:00 +0000 UTC Processing...
	// 0001-01-01 00:00:01 +0000 UTC Processing...
	// 0001-01-01 00:00:02 +0000 UTC Processing...
	// 0001-01-01 00:00:03 +0000 UTC Processing...
	// 0001-01-01 00:00:04 +0000 UTC Processing...
	// 0001-01-01 00:00:05 +0000 UTC Processing...
	// 0001-01-01 00:00:06 +0000 UTC Processing...
	// 0001-01-01 00:00:07 +0000 UTC Processing...
	// 0001-01-01 00:00:08 +0000 UTC Processing...
	// 0001-01-01 00:00:09 +0000 UTC Processing...
}

// ExampleBinarySemaphore demonstrates how to use the BinarySemaphore to synchronize actions. It simulates processing ten (10) items, one (1) at a time.
func ExampleBinarySemaphore() {
	sim := NewSimulation()
	sem := NewBinarySemaphore(sim)

	// Simulate processing ten (10) items, one (1) at a time.
	timeToProcess := 10 * time.Second
	for i := range 10 {
		sem.Acquire(func(sim *Simulation) {
			// We have now acquired the semaphore and can start processing.
			fmt.Println(sim.Now, "Processing item", i)

			sim.Schedule(Event{When: sim.Now.Add(timeToProcess), Action: func(sim *Simulation) {
				fmt.Println(sim.Now, "Done processing item", i)
				sem.Release()
			}})
		})
	}
	sim.RunUntilDone()

	// Output:
	// 0001-01-01 00:00:00 +0000 UTC Processing item 0
	// 0001-01-01 00:00:10 +0000 UTC Done processing item 0
	// 0001-01-01 00:00:10 +0000 UTC Processing item 1
	// 0001-01-01 00:00:20 +0000 UTC Done processing item 1
	// 0001-01-01 00:00:20 +0000 UTC Processing item 2
	// 0001-01-01 00:00:30 +0000 UTC Done processing item 2
	// 0001-01-01 00:00:30 +0000 UTC Processing item 3
	// 0001-01-01 00:00:40 +0000 UTC Done processing item 3
	// 0001-01-01 00:00:40 +0000 UTC Processing item 4
	// 0001-01-01 00:00:50 +0000 UTC Done processing item 4
	// 0001-01-01 00:00:50 +0000 UTC Processing item 5
	// 0001-01-01 00:01:00 +0000 UTC Done processing item 5
	// 0001-01-01 00:01:00 +0000 UTC Processing item 6
	// 0001-01-01 00:01:10 +0000 UTC Done processing item 6
	// 0001-01-01 00:01:10 +0000 UTC Processing item 7
	// 0001-01-01 00:01:20 +0000 UTC Done processing item 7
	// 0001-01-01 00:01:20 +0000 UTC Processing item 8
	// 0001-01-01 00:01:30 +0000 UTC Done processing item 8
	// 0001-01-01 00:01:30 +0000 UTC Processing item 9
	// 0001-01-01 00:01:40 +0000 UTC Done processing item 9
}

// ExampleCountingSemaphore demonstrates how to use the CountingSemaphore to synchronize actions.
func ExampleCountingSemaphore() {
	sim := NewSimulation()
	sem := NewCountingSemaphore(sim, 3)

	// Simulate processing 10 items, 3 at a time.
	timeToProcess := 10 * time.Second
	for i := range 10 {
		sem.Acquire(func(sim *Simulation) {
			// We have now acquired the semaphore and can start processing.
			fmt.Println(sim.Now, "Processing item", i)

			sim.Schedule(Event{When: sim.Now.Add(timeToProcess), Action: func(sim *Simulation) {
				fmt.Println(sim.Now, "Done processing item", i)
				sem.Release()
			}})
		})
	}
	sim.RunUntilDone()

	// Output:
	// 0001-01-01 00:00:00 +0000 UTC Processing item 0
	// 0001-01-01 00:00:00 +0000 UTC Processing item 1
	// 0001-01-01 00:00:00 +0000 UTC Processing item 2
	// 0001-01-01 00:00:10 +0000 UTC Done processing item 0
	// 0001-01-01 00:00:10 +0000 UTC Processing item 3
	// 0001-01-01 00:00:10 +0000 UTC Done processing item 1
	// 0001-01-01 00:00:10 +0000 UTC Processing item 4
	// 0001-01-01 00:00:10 +0000 UTC Done processing item 2
	// 0001-01-01 00:00:10 +0000 UTC Processing item 5
	// 0001-01-01 00:00:20 +0000 UTC Done processing item 3
	// 0001-01-01 00:00:20 +0000 UTC Processing item 6
	// 0001-01-01 00:00:20 +0000 UTC Done processing item 4
	// 0001-01-01 00:00:20 +0000 UTC Processing item 7
	// 0001-01-01 00:00:20 +0000 UTC Done processing item 5
	// 0001-01-01 00:00:20 +0000 UTC Processing item 8
	// 0001-01-01 00:00:30 +0000 UTC Done processing item 6
	// 0001-01-01 00:00:30 +0000 UTC Processing item 9
	// 0001-01-01 00:00:30 +0000 UTC Done processing item 7
	// 0001-01-01 00:00:30 +0000 UTC Done processing item 8
	// 0001-01-01 00:00:40 +0000 UTC Done processing item 9
}
