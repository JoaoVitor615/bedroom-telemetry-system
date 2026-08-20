package model

import "sync"

type CurrentState struct {
	mu          sync.RWMutex
	Temperature string `json:"temperature"`
	DoorState   string `json:"door_state"`
}

func (s *CurrentState) UpdateTemp(temp string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Temperature = temp
}

func (s *CurrentState) UpdateDoor(state string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.DoorState = state
}
