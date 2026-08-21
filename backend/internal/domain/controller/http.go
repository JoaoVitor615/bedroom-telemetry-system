package controller

import (
	"encoding/json"
	"net/http"

	"github.com/JoaoVitor615/bedroom-telemetry-system/internal/domain"
)

type HTTPController struct {
	service domain.IBedroomService
}

// NewHTTPController initializes the controller with the service dependency
func NewHTTPController(service domain.IBedroomService) *HTTPController {
	return &HTTPController{service: service}
}

// GetCurrentState handles requests for instantaneous telemetry
func (c *HTTPController) GetCurrentState(w http.ResponseWriter, r *http.Request) {
	response := c.service.GetCurrentState()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetHistory handles requests for the door events log
func (c *HTTPController) GetHistory(w http.ResponseWriter, r *http.Request) {
	response, err := c.service.GetDoorHistory(r.Context(), 20)
	if err != nil {
		http.Error(w, "Failed to fetch history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
