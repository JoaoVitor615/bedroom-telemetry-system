package service

import (
	"context"
	"time"

	"github.com/JoaoVitor615/bedroom-telemetry-system/internal/domain"
	"github.com/JoaoVitor615/bedroom-telemetry-system/internal/domain/repository"
	"github.com/JoaoVitor615/bedroom-telemetry-system/internal/dto"
)

type BedroomService struct {
	repo  domain.IDoorRepository
	state *repository.CurrentState
}

// NewBedroomService creates a new instance of the business logic service
func NewBedroomService(repo domain.IDoorRepository, state *repository.CurrentState) domain.IBedroomService {
	return &BedroomService{
		repo:  repo,
		state: state,
	}
}

func (s *BedroomService) ProcessTemperature(temp string) {
	s.state.SetTemperature(temp)
}

func (s *BedroomService) ProcessDoorEvent(ctx context.Context, state string) error {
	// Update real-time state
	s.state.SetDoorState(state)

	// Create domain entity
	event := repository.DoorEvent{
		State:     state,
		Timestamp: time.Now(),
	}

	// Persist using the abstracted repository interface
	return s.repo.SaveEvent(ctx, event)
}

func (s *BedroomService) GetCurrentState() dto.StateResponseDTO {
	temp, door := s.state.Get()
	return dto.StateResponseDTO{
		Temperature: temp,
		DoorState:   door,
	}
}

func (s *BedroomService) GetDoorHistory(ctx context.Context, limit int64) ([]dto.DoorEventDTO, error) {
	events, err := s.repo.GetHistory(ctx, limit)
	if err != nil {
		return nil, err
	}

	// Map Domain Entities to DTOs
	var dtos []dto.DoorEventDTO
	for _, e := range events {
		dtos = append(dtos, dto.DoorEventDTO{
			State:     e.State,
			Timestamp: e.Timestamp,
		})
	}

	if dtos == nil {
		dtos = make([]dto.DoorEventDTO, 0) // Return empty array instead of null
	}

	return dtos, nil
}
