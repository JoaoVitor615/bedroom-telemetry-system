package repository

import (
	"fmt"
	"sync"
	"testing"
)

func TestCurrentState_SetAndGet(t *testing.T) {
	state := &CurrentState{}

	// Verify initial zero-values
	temp, door := state.Get()
	if temp != "" || door != "" {
		t.Errorf("expected empty initial state, got temp=%q, door=%q", temp, door)
	}

	// Update temperature
	state.SetTemperature("24.5")
	temp, door = state.Get()
	if temp != "24.5" || door != "" {
		t.Errorf("expected temp=24.5, door='', got temp=%q, door=%q", temp, door)
	}

	// Update door state
	state.SetDoorState("OPENED")
	temp, door = state.Get()
	if temp != "24.5" || door != "OPENED" {
		t.Errorf("expected temp=24.5, door=OPENED, got temp=%q, door=%q", temp, door)
	}
}

func TestCurrentState_ConcurrentAccess(t *testing.T) {
	state := &CurrentState{}
	var wg sync.WaitGroup

	numGoroutines := 50
	for i := 0; i < numGoroutines; i++ {
		wg.Add(3)

		go func(val int) {
			defer wg.Done()
			state.SetTemperature(fmt.Sprintf("%d.0", val))
		}(i)

		go func(val int) {
			defer wg.Done()
			if val%2 == 0 {
				state.SetDoorState("OPENED")
			} else {
				state.SetDoorState("CLOSED")
			}
		}(i)

		go func() {
			defer wg.Done()
			_, _ = state.Get()
		}()
	}

	wg.Wait()

	// Ensure final read doesn't panic or data race
	finalTemp, finalDoor := state.Get()
	if finalTemp == "" || finalDoor == "" {
		t.Errorf("expected non-empty state after concurrent updates, got temp=%q, door=%q", finalTemp, finalDoor)
	}
}
