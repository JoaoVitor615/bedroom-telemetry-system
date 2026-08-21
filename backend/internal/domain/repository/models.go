package repository

import (
	"sync"
	"time"
)

// DoorEvent represents the core entity for a door state change
type DoorEvent struct {
	State     string
	Timestamp time.Time
}

// CurrentState holds the thread-safe instantaneous state of the bedroom
type CurrentState struct {
	mu          sync.RWMutex
	Temperature string
	DoorState   string
}

// SetTemperature safely updates the temperature
func (s *CurrentState) SetTemperature(temp string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Temperature = temp
}

// SetDoorState safely updates the door state
func (s *CurrentState) SetDoorState(state string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.DoorState = state
}

// Get safely retrieves both values
func (s *CurrentState) Get() (string, string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Temperature, s.DoorState
}
