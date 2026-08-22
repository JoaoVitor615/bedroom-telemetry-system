package controller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/JoaoVitor615/bedroom-telemetry-system/internal/dto"
)

// mockBedroomService implements domain.IBedroomService for testing HTTP responses
type mockBedroomService struct {
	processTemperatureFunc func(temp string)
	processDoorEventFunc   func(ctx context.Context, state string) error
	getCurrentStateFunc    func() dto.StateResponseDTO
	getDoorHistoryFunc     func(ctx context.Context, limit int64) ([]dto.DoorEventDTO, error)
}

func (m *mockBedroomService) ProcessTemperature(temp string) {
	if m.processTemperatureFunc != nil {
		m.processTemperatureFunc(temp)
	}
}

func (m *mockBedroomService) ProcessDoorEvent(ctx context.Context, state string) error {
	if m.processDoorEventFunc != nil {
		return m.processDoorEventFunc(ctx, state)
	}
	return nil
}

func (m *mockBedroomService) GetCurrentState() dto.StateResponseDTO {
	if m.getCurrentStateFunc != nil {
		return m.getCurrentStateFunc()
	}
	return dto.StateResponseDTO{}
}

func (m *mockBedroomService) GetDoorHistory(ctx context.Context, limit int64) ([]dto.DoorEventDTO, error) {
	if m.getDoorHistoryFunc != nil {
		return m.getDoorHistoryFunc(ctx, limit)
	}
	return nil, nil
}

func TestHTTPController_GetCurrentState(t *testing.T) {
	mockSvc := &mockBedroomService{
		getCurrentStateFunc: func() dto.StateResponseDTO {
			return dto.StateResponseDTO{
				Temperature: "22.5",
				DoorState:   "CLOSED",
			}
		},
	}

	ctrl := NewHTTPController(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/state", nil)
	rr := httptest.NewRecorder()

	ctrl.GetCurrentState(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", contentType)
	}

	var resp dto.StateResponseDTO
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	if resp.Temperature != "22.5" || resp.DoorState != "CLOSED" {
		t.Errorf("expected temp=22.5, door=CLOSED, got temp=%q, door=%q", resp.Temperature, resp.DoorState)
	}
}

func TestHTTPController_GetHistory_Success(t *testing.T) {
	fixedTime := time.Date(2026, 8, 22, 10, 0, 0, 0, time.UTC)
	mockSvc := &mockBedroomService{
		getDoorHistoryFunc: func(ctx context.Context, limit int64) ([]dto.DoorEventDTO, error) {
			if limit != 20 {
				t.Errorf("expected limit 20, got %d", limit)
			}
			return []dto.DoorEventDTO{
				{State: "OPENED", Timestamp: fixedTime},
			}, nil
		},
	}

	ctrl := NewHTTPController(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/history", nil)
	rr := httptest.NewRecorder()

	ctrl.GetHistory(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %d", rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", contentType)
	}

	var resp []dto.DoorEventDTO
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	if len(resp) != 1 {
		t.Fatalf("expected 1 element, got %d", len(resp))
	}
	if resp[0].State != "OPENED" || !resp[0].Timestamp.Equal(fixedTime) {
		t.Errorf("unexpected event response: %+v", resp[0])
	}
}

func TestHTTPController_GetHistory_Error(t *testing.T) {
	mockSvc := &mockBedroomService{
		getDoorHistoryFunc: func(ctx context.Context, limit int64) ([]dto.DoorEventDTO, error) {
			return nil, errors.New("db error")
		},
	}

	ctrl := NewHTTPController(mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/api/history", nil)
	rr := httptest.NewRecorder()

	ctrl.GetHistory(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected status code 500, got %d", rr.Code)
	}
}
