package handlers

import (
	"encoding/json"
	"net/http"

	"city/internal/models"
	"city/internal/services"
	"city/internal/storage"
)

type ResourcesHandler struct {
	service *services.ResourcesService
}

func NewResourcesHandler(store storage.Store) *ResourcesHandler {
	return &ResourcesHandler{
		service: services.NewResourcesService(store),
	}
}

// Trade godoc
// @Summary Trade resources
// @Description Trade resources (buy/sell) with specified price
// @Tags resources
// @Accept json
// @Produce json
// @Param request body models.TradeRequest true "Trade request"
// @Success 200 {object} models.Resources
// @Failure 400 {object} models.ErrorResponse
// @Router /city/trade [post]
func (h *ResourcesHandler) Trade(w http.ResponseWriter, r *http.Request) {
	var req models.TradeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resources, err := h.service.Trade(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resources)
}

// History godoc
// @Summary Get resource history
// @Description Get history of resource changes over time
// @Tags resources
// @Produce json
// @Success 200 {array} models.ResourceSnapshot
// @Router /city/resources/history [get]
func (h *ResourcesHandler) History(w http.ResponseWriter, r *http.Request) {
	history, err := h.service.GetResourceHistory()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

// Adjust godoc
// @Summary Adjust resources
// @Description Manually adjust resource amounts
// @Tags resources
// @Accept json
// @Produce json
// @Param request body models.AdjustResourcesRequest true "Resource adjustment request"
// @Success 200 {object} models.Resources
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /city/resources/adjust [patch]
func (h *ResourcesHandler) Adjust(w http.ResponseWriter, r *http.Request) {
	var req models.AdjustResourcesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resources, err := h.service.AdjustResources(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resources)
}
