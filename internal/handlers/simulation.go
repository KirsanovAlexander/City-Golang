package handlers

import (
	"encoding/json"
	"net/http"

	"city/internal/models"
	"city/internal/services"
	"city/internal/storage"
)

type SimulationHandler struct {
	service *services.SimulationService
}

func NewSimulationHandler(store storage.Store) *SimulationHandler {
	return &SimulationHandler{
		service: services.NewSimulationService(store),
	}
}

// Tick godoc
// @Summary Advance simulation by one day
// @Description Process one day of simulation (production, consumption, events)
// @Tags simulation
// @Produce json
// @Success 200 {object} models.City
// @Failure 404 {object} models.ErrorResponse
// @Router /city/tick [post]
func (h *SimulationHandler) Tick(w http.ResponseWriter, r *http.Request) {
	city, err := h.service.Tick()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(city)
}

// RandomEvent godoc
// @Summary Trigger random event
// @Description Trigger a random event that affects the city
// @Tags simulation
// @Produce json
// @Success 200 {object} models.Event
// @Failure 404 {object} models.ErrorResponse
// @Router /city/events/random [post]
func (h *SimulationHandler) RandomEvent(w http.ResponseWriter, r *http.Request) {
	event, err := h.service.RandomEvent()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

// CustomEvent godoc
// @Summary Trigger custom event
// @Description Trigger a custom event with specified parameters
// @Tags simulation
// @Accept json
// @Produce json
// @Param request body models.CustomEventRequest true "Custom event request"
// @Success 200 {object} models.Event
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /city/events/custom [post]
func (h *SimulationHandler) CustomEvent(w http.ResponseWriter, r *http.Request) {
	var req models.CustomEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := h.service.CustomEvent(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

// EventsHistory godoc
// @Summary Get events history
// @Description Get history of all events that occurred
// @Tags simulation
// @Produce json
// @Success 200 {array} models.Event
// @Router /city/events/history [get]
func (h *SimulationHandler) EventsHistory(w http.ResponseWriter, r *http.Request) {
	events, err := h.service.GetEventsHistory()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// Stats godoc
// @Summary Get city statistics
// @Description Get current city statistics and metrics
// @Tags simulation
// @Produce json
// @Success 200 {object} models.Stats
// @Failure 404 {object} models.ErrorResponse
// @Router /city/stats [get]
func (h *SimulationHandler) Stats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
