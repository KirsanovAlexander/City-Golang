package handlers

import (
	"encoding/json"
	"net/http"

	"city/internal/models"
	"city/internal/services"
	"city/internal/storage"
)

type CityHandler struct {
	service *services.CityService
}

func NewCityHandler(store storage.Store) *CityHandler {
	return &CityHandler{
		service: services.NewCityService(store),
	}
}

// Create godoc
// @Summary Create a new city
// @Description Create a new city with specified name and difficulty
// @Tags city
// @Accept json
// @Produce json
// @Param request body models.CreateCityRequest true "City creation request"
// @Success 201 {object} models.City
// @Failure 400 {object} models.ErrorResponse
// @Router /city [post]
func (h *CityHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	city, err := h.service.CreateCity(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(city)
}

// Get godoc
// @Summary Get city information
// @Description Get current city state and information
// @Tags city
// @Produce json
// @Success 200 {object} models.City
// @Failure 404 {object} models.ErrorResponse
// @Router /city [get]
func (h *CityHandler) Get(w http.ResponseWriter, r *http.Request) {
	city, err := h.service.GetCity()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(city)
}

// Reset godoc
// @Summary Reset city
// @Description Reset the city to initial state
// @Tags city
// @Success 204
// @Router /city/reset [delete]
func (h *CityHandler) Reset(w http.ResponseWriter, r *http.Request) {
	err := h.service.ResetCity()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// UpdateSettings godoc
// @Summary Update city settings
// @Description Update city configuration settings
// @Tags city
// @Accept json
// @Produce json
// @Param request body models.UpdateSettingsRequest true "Settings update request"
// @Success 200 {object} models.CitySettings
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /city/settings [patch]
func (h *CityHandler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	settings, err := h.service.UpdateSettings(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}
