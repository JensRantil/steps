package events

import (
	"testing"
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

// ExampleCondition demonstrates how to use the Condition to synchronize actions.
func ExampleCondition() {
	// TODO: Implement this.
}
