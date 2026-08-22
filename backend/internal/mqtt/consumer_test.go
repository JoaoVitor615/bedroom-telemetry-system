package mqtt

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/JoaoVitor615/bedroom-telemetry-system/internal/dto"
)

type mockBedroomService struct {
	mu                     sync.Mutex
	lastTemp               string
	lastDoor               string
	processTemperatureDone chan struct{}
	processDoorEventDone   chan struct{}
}

func (m *mockBedroomService) ProcessTemperature(temp string) {
	m.mu.Lock()
	m.lastTemp = temp
	m.mu.Unlock()
	if m.processTemperatureDone != nil {
		select {
		case m.processTemperatureDone <- struct{}{}:
		default:
		}
	}
}

func (m *mockBedroomService) ProcessDoorEvent(ctx context.Context, state string) error {
	m.mu.Lock()
	m.lastDoor = state
	m.mu.Unlock()
	if m.processDoorEventDone != nil {
		select {
		case m.processDoorEventDone <- struct{}{}:
		default:
		}
	}
	return nil
}

func (m *mockBedroomService) GetCurrentState() dto.StateResponseDTO {
	m.mu.Lock()
	defer m.mu.Unlock()
	return dto.StateResponseDTO{
		Temperature: m.lastTemp,
		DoorState:   m.lastDoor,
	}
}

func (m *mockBedroomService) GetDoorHistory(ctx context.Context, limit int64) ([]dto.DoorEventDTO, error) {
	return nil, nil
}

func TestConsumer_ServiceIntegration(t *testing.T) {
	mockSvc := &mockBedroomService{
		processTemperatureDone: make(chan struct{}, 1),
		processDoorEventDone:   make(chan struct{}, 1),
	}

	// Test temperature processing
	mockSvc.ProcessTemperature("27.5")
	state := mockSvc.GetCurrentState()
	if state.Temperature != "27.5" {
		t.Errorf("expected temperature 27.5, got %q", state.Temperature)
	}

	// Test door event processing
	err := mockSvc.ProcessDoorEvent(context.Background(), "OPENED")
	if err != nil {
		t.Fatalf("expected no error processing door event, got %v", err)
	}

	// Wait briefly for processing
	time.Sleep(10 * time.Millisecond)

	state = mockSvc.GetCurrentState()
	if state.DoorState != "OPENED" {
		t.Errorf("expected door state OPENED, got %q", state.DoorState)
	}
}
