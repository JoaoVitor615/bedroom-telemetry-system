package model

import (
	"sync"
	"time"
)

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

func (s *CurrentState) Get() (string, string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Temperature, s.DoorState
}

type DoorEvent struct {
	State     string    `bson:"state" json:"state"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}
