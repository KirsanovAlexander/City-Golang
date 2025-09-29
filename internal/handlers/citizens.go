package handlers

import (
	"encoding/json"
	"net/http"

	"city/internal/models"
	"city/internal/services"
	"city/internal/storage"

	"github.com/go-chi/chi/v5"
)

type CitizensHandler struct {
	service *services.CitizensService
}

func NewCitizensHandler(store storage.Store) *CitizensHandler {
	return &CitizensHandler{
		service: services.NewCitizensService(store),
	}
}

// Create godoc
// @Summary Create a new citizen
// @Description Create a new citizen with specified name and job
// @Tags citizens
// @Accept json
// @Produce json
// @Param request body models.CitizenCreateRequest true "Citizen creation request"
// @Success 201 {object} models.Citizen
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /city/citizens [post]
func (h *CitizensHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CitizenCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	citizen, err := h.service.CreateCitizen(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(citizen)
}

// List godoc
// @Summary List all citizens
// @Description Get list of all citizens in the city
// @Tags citizens
// @Produce json
// @Success 200 {array} models.Citizen
// @Failure 404 {object} models.ErrorResponse
// @Router /city/citizens [get]
func (h *CitizensHandler) List(w http.ResponseWriter, r *http.Request) {
	citizens, err := h.service.ListCitizens()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(citizens)
}

// ChangeJob godoc
// @Summary Change citizen job
// @Description Change a citizen's job by ID
// @Tags citizens
// @Accept json
// @Produce json
// @Param id path string true "Citizen ID"
// @Param request body models.ChangeJobRequest true "Job change request"
// @Success 200 {object} models.Citizen
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /city/citizens/{id}/job [patch]
func (h *CitizensHandler) ChangeJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.ChangeJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	citizen, err := h.service.ChangeJob(id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(citizen)
}

// ChangeHappiness godoc
// @Summary Change citizen happiness
// @Description Adjust a citizen's happiness by ID
// @Tags citizens
// @Accept json
// @Produce json
// @Param id path string true "Citizen ID"
// @Param request body models.ChangeHappinessRequest true "Happiness change request"
// @Success 200 {object} models.Citizen
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /city/citizens/{id}/happiness [patch]
func (h *CitizensHandler) ChangeHappiness(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.ChangeHappinessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	citizen, err := h.service.ChangeHappiness(id, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(citizen)
}

// Delete godoc
// @Summary Delete a citizen
// @Description Delete a citizen by ID
// @Tags citizens
// @Param id path string true "Citizen ID"
// @Success 204
// @Failure 404 {object} models.ErrorResponse
// @Router /city/citizens/{id} [delete]
func (h *CitizensHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.service.DeleteCitizen(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// JobsStats godoc
// @Summary Get job statistics
// @Description Get statistics of citizens by job type
// @Tags citizens
// @Produce json
// @Success 200 {object} models.JobsStatsResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /city/citizens/jobs [get]
func (h *CitizensHandler) JobsStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.service.GetJobsStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// MassAdd godoc
// @Summary Mass add citizens
// @Description Add multiple citizens with specified job and prefix
// @Tags citizens
// @Accept json
// @Produce json
// @Param request body models.MassAddRequest true "Mass add request"
// @Success 201 {array} models.Citizen
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /city/citizens/mass-add [post]
func (h *CitizensHandler) MassAdd(w http.ResponseWriter, r *http.Request) {
	var req models.MassAddRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	citizens, err := h.service.MassAddCitizens(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(citizens)
}
