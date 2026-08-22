package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/JoaoVitor615/bedroom-telemetry-system/internal/domain/repository"
	"github.com/JoaoVitor615/bedroom-telemetry-system/internal/domain/service"
)

// mockDoorRepository implements domain.IDoorRepository for testing
type mockDoorRepository struct {
	saveEventFunc  func(ctx context.Context, event repository.DoorEvent) error
	getHistoryFunc func(ctx context.Context, limit int64) ([]repository.DoorEvent, error)
}

func (m *mockDoorRepository) SaveEvent(ctx context.Context, event repository.DoorEvent) error {
	if m.saveEventFunc != nil {
		return m.saveEventFunc(ctx, event)
	}
	return nil
}

func (m *mockDoorRepository) GetHistory(ctx context.Context, limit int64) ([]repository.DoorEvent, error) {
	if m.getHistoryFunc != nil {
		return m.getHistoryFunc(ctx, limit)
	}
	return nil, nil
}

func TestBedroomService_ProcessTemperature(t *testing.T) {
	repo := &mockDoorRepository{}
	state := &repository.CurrentState{}
	svc := service.NewBedroomService(repo, state)

	svc.ProcessTemperature("26.8")

	dto := svc.GetCurrentState()
	if dto.Temperature != "26.8" {
		t.Errorf("expected temperature 26.8, got %q", dto.Temperature)
	}
}

func TestBedroomService_ProcessDoorEvent_Success(t *testing.T) {
	var savedEvent repository.DoorEvent
	repo := &mockDoorRepository{
		saveEventFunc: func(ctx context.Context, event repository.DoorEvent) error {
			savedEvent = event
			return nil
		},
	}
	state := &repository.CurrentState{}
	svc := service.NewBedroomService(repo, state)

	err := svc.ProcessDoorEvent(context.Background(), "OPENED")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	dto := svc.GetCurrentState()
	if dto.DoorState != "OPENED" {
		t.Errorf("expected door state OPENED, got %q", dto.DoorState)
	}

	if savedEvent.State != "OPENED" {
		t.Errorf("expected saved event state OPENED, got %q", savedEvent.State)
	}
	if savedEvent.Timestamp.IsZero() {
		t.Errorf("expected non-zero timestamp in saved event")
	}
}

func TestBedroomService_ProcessDoorEvent_RepoError(t *testing.T) {
	expectedErr := errors.New("database unavailable")
	repo := &mockDoorRepository{
		saveEventFunc: func(ctx context.Context, event repository.DoorEvent) error {
			return expectedErr
		},
	}
	state := &repository.CurrentState{}
	svc := service.NewBedroomService(repo, state)

	err := svc.ProcessDoorEvent(context.Background(), "CLOSED")
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	// State should still be updated in real-time memory even if DB save fails
	dto := svc.GetCurrentState()
	if dto.DoorState != "CLOSED" {
		t.Errorf("expected door state CLOSED in memory, got %q", dto.DoorState)
	}
}

func TestBedroomService_GetCurrentState(t *testing.T) {
	repo := &mockDoorRepository{}
	state := &repository.CurrentState{}
	state.SetTemperature("23.0")
	state.SetDoorState("CLOSED")

	svc := service.NewBedroomService(repo, state)
	dto := svc.GetCurrentState()

	if dto.Temperature != "23.0" || dto.DoorState != "CLOSED" {
		t.Errorf("expected temp=23.0, door=CLOSED, got temp=%q, door=%q", dto.Temperature, dto.DoorState)
	}
}

func TestBedroomService_GetDoorHistory_Success(t *testing.T) {
	now := time.Now()
	repoEvents := []repository.DoorEvent{
		{State: "OPENED", Timestamp: now},
		{State: "CLOSED", Timestamp: now.Add(-5 * time.Minute)},
	}

	repo := &mockDoorRepository{
		getHistoryFunc: func(ctx context.Context, limit int64) ([]repository.DoorEvent, error) {
			if limit != 10 {
				t.Errorf("expected limit 10, got %d", limit)
			}
			return repoEvents, nil
		},
	}
	state := &repository.CurrentState{}
	svc := service.NewBedroomService(repo, state)

	dtos, err := svc.GetDoorHistory(context.Background(), 10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(dtos) != 2 {
		t.Fatalf("expected 2 DTOs, got %d", len(dtos))
	}
	if dtos[0].State != "OPENED" || !dtos[0].Timestamp.Equal(now) {
		t.Errorf("unexpected first DTO: %+v", dtos[0])
	}
	if dtos[1].State != "CLOSED" {
		t.Errorf("unexpected second DTO: %+v", dtos[1])
	}
}

func TestBedroomService_GetDoorHistory_EmptyList(t *testing.T) {
	repo := &mockDoorRepository{
		getHistoryFunc: func(ctx context.Context, limit int64) ([]repository.DoorEvent, error) {
			return nil, nil
		},
	}
	state := &repository.CurrentState{}
	svc := service.NewBedroomService(repo, state)

	dtos, err := svc.GetDoorHistory(context.Background(), 20)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if dtos == nil {
		t.Fatal("expected non-null slice for empty history, got nil")
	}
	if len(dtos) != 0 {
		t.Errorf("expected empty slice, got length %d", len(dtos))
	}
}

func TestBedroomService_GetDoorHistory_RepoError(t *testing.T) {
	expectedErr := errors.New("query failed")
	repo := &mockDoorRepository{
		getHistoryFunc: func(ctx context.Context, limit int64) ([]repository.DoorEvent, error) {
			return nil, expectedErr
		},
	}
	state := &repository.CurrentState{}
	svc := service.NewBedroomService(repo, state)

	dtos, err := svc.GetDoorHistory(context.Background(), 20)
	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
	if dtos != nil {
		t.Errorf("expected nil DTOs on error, got %v", dtos)
	}
}
