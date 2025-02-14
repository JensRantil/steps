package events

import (
	"container/heap"
	"fmt"
	"time"
)

type ConditionActionID int

// Condition is a condition that can be used to synchronize actions. Multiple actions can be waiting for the same condition, and be triggered by a Signal or Broadcast, similarly to sync.Cond.
type Condition struct {
	sim    *Simulation
	heap   *conditionHeap
	nextID ConditionActionID
}

// NewCondition creates a new condition.
func NewCondition(sim *Simulation) *Condition {
	return &Condition{
		sim:  sim,
		heap: newConditionHeap(),
	}
}

func (c *Condition) Wait(a Action) ConditionActionID {
	id := c.nextID
	heap.Push(c.heap, conditionActionItem{ID: id, Action: a})
	c.nextID++
	return id
}

// Cancel cancels an action waiting for this condition. Returns true if the action was found and removed, false otherwise (e.g. it was already executed, never existed, or was previously cancelled).
func (c *Condition) Cancel(id ConditionActionID) bool {
	index, found := c.heap.IndexByID[id]
	if !found {
		return false
	}
	heap.Remove(c.heap, index)
	return true
}

// Signal wakes up one action waiting for this condition. Actions are woken up in the order they were waiting.
func (c *Condition) Signal() {
	if c.heap.Len() == 0 {
		return
	}

	item := heap.Pop(c.heap).(conditionActionItem)
	c.sim.Schedule(ScheduledEvent{
		When:   time.Time{}, // As soon as possible.
		Action: item.Action,
	})
}

// Broadcast wakes up all actions waiting for this condition. Actions are woken up in the order they were waiting.
func (c *Condition) Broadcast() {
	for c.heap.Len() > 0 {
		// It's important that we iterate over the heap in order to schedule in a FIFO manner.
		item := heap.Pop(c.heap).(conditionActionItem)
		c.sim.Schedule(ScheduledEvent{
			When:   time.Time{}, // As soon as possible.
			Action: item.Action,
		})
	}
}

// conditionActionItem is an item in the condition heap.
type conditionActionItem struct {
	ID     ConditionActionID
	Action Action
}

// A heap of actions, ordered by ConditionActionID.
type conditionHeap struct {
	items     []conditionActionItem
	IndexByID map[ConditionActionID]int
}

func newConditionHeap() *conditionHeap {
	return &conditionHeap{
		IndexByID: make(map[ConditionActionID]int),
	}
}

func (h conditionHeap) Len() int {
	sliceLen := len(h.items)
	if sliceLen != len(h.IndexByID) {
		panic(fmt.Sprintf("len(h.Events) != len(h.IndexByID): %d != %d", sliceLen, len(h.IndexByID)))
	}
	return sliceLen
}
func (h conditionHeap) Less(i, j int) bool {
	return h.items[i].ID < h.items[j].ID
}
func (h conditionHeap) Swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]

	h.IndexByID[h.items[i].ID] = i
	h.IndexByID[h.items[j].ID] = j
}

func (h *conditionHeap) Push(xUntyped any) {
	x := xUntyped.(conditionActionItem)
	if _, found := h.IndexByID[x.ID]; found {
		panic(fmt.Sprintf("event with ID %d already exists", x.ID))
	}
	h.items = append(h.items, x)
	h.IndexByID[x.ID] = len(h.items) - 1
}

func (h *conditionHeap) Pop() any {
	n := len(h.items)
	x := h.items[n-1]
	h.items = h.items[0 : n-1]
	delete(h.IndexByID, x.ID)
	return x
}
