package dto

import "time"

// StateResponseDTO defines the JSON structure for the current state endpoint
type StateResponseDTO struct {
	Temperature string `json:"temperature"`
	DoorState   string `json:"door_state"`
}

// DoorEventDTO defines the JSON structure for the history endpoint
type DoorEventDTO struct {
	State     string    `json:"state"`
	Timestamp time.Time `json:"timestamp"`
}
