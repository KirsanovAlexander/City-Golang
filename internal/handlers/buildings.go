package handlers

import (
	"encoding/json"
	"net/http"

	"city/internal/models"
	"city/internal/services"
	"city/internal/storage"

	"github.com/go-chi/chi/v5"
)

type BuildingsHandler struct {
	service *services.BuildingsService
}

func NewBuildingsHandler(store storage.Store) *BuildingsHandler {
	return &BuildingsHandler{
		service: services.NewBuildingsService(store),
	}
}

// Create godoc
// @Summary Create a new building
// @Description Create a new building of specified type
// @Tags buildings
// @Accept json
// @Produce json
// @Param request body models.BuildRequest true "Building creation request"
// @Success 201 {object} models.Building
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /city/buildings [post]
func (h *BuildingsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.BuildRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	building, err := h.service.CreateBuilding(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(building)
}

// Upgrade godoc
// @Summary Upgrade a building
// @Description Upgrade a building by ID to next level
// @Tags buildings
// @Produce json
// @Param id path string true "Building ID"
// @Success 200 {object} models.Building
// @Failure 404 {object} models.ErrorResponse
// @Router /city/buildings/{id}/upgrade [patch]
func (h *BuildingsHandler) Upgrade(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	building, err := h.service.UpgradeBuilding(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(building)
}

// Repair godoc
// @Summary Repair a building
// @Description Repair a building by ID with specified amount
// @Tags buildings
// @Accept json
// @Produce json
// @Param id path string true "Building ID"
// @Param request body models.RepairRequest true "Repair request"
// @Success 200 {object} models.Building
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /city/buildings/{id}/repair [patch]
func (h *BuildingsHandler) Repair(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.RepairRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	building, err := h.service.RepairBuilding(id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(building)
}

// Delete godoc
// @Summary Delete a building
// @Description Delete a building by ID
// @Tags buildings
// @Param id path string true "Building ID"
// @Success 204
// @Failure 404 {object} models.ErrorResponse
// @Router /city/buildings/{id} [delete]
func (h *BuildingsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.DeleteBuilding(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// List godoc
// @Summary List all buildings
// @Description Get list of all buildings in the city
// @Tags buildings
// @Produce json
// @Success 200 {array} models.Building
// @Failure 404 {object} models.ErrorResponse
// @Router /city/buildings [get]
func (h *BuildingsHandler) List(w http.ResponseWriter, r *http.Request) {
	buildings, err := h.service.ListBuildings()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(buildings)
}

// Effects godoc
// @Summary Get building effects
// @Description Get description of effects for each building type
// @Tags buildings
// @Produce json
// @Success 200 {array} models.BuildingEffect
// @Router /city/buildings/effects [get]
func (h *BuildingsHandler) Effects(w http.ResponseWriter, r *http.Request) {
	effects := h.service.GetBuildingEffects()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(effects)
}
