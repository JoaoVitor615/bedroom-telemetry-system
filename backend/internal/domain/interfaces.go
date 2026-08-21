package domain

import (
	"context"

	"github.com/JoaoVitor615/bedroom-telemetry-system/internal/domain/repository"
	"github.com/JoaoVitor615/bedroom-telemetry-system/internal/dto"
)

// IDoorRepository defines the contract for data persistence (Database agnostic)
type IDoorRepository interface {
	SaveEvent(ctx context.Context, event repository.DoorEvent) error
	GetHistory(ctx context.Context, limit int64) ([]repository.DoorEvent, error)
}

// IBedroomService defines the contract for business logic and use cases
type IBedroomService interface {
	ProcessTemperature(temp string)
	ProcessDoorEvent(ctx context.Context, state string) error
	GetCurrentState() dto.StateResponseDTO
	GetDoorHistory(ctx context.Context, limit int64) ([]dto.DoorEventDTO, error)
}
