package events

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
