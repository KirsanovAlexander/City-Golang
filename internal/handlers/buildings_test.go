package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"city/internal/models"
	"city/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildingsHandler_Create(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewBuildingsHandler(store)

	// Create a city first
	store.CreateCity("Test City", "normal")

	tests := []struct {
		name           string
		request        models.BuildRequest
		expectedStatus int
	}{
		{
			name: "create farm building",
			request: models.BuildRequest{
				Type: models.BuildingFarm,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "create house building",
			request: models.BuildRequest{
				Type: models.BuildingHouse,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "create factory building",
			request: models.BuildRequest{
				Type: models.BuildingFactory,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "create power plant building",
			request: models.BuildRequest{
				Type: models.BuildingPowerPlant,
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/city/buildings", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Create(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusCreated {
				var building models.Building
				err := json.Unmarshal(w.Body.Bytes(), &building)
				require.NoError(t, err)
				assert.Equal(t, tt.request.Type, building.Type)
				assert.Equal(t, 1, building.Level)
				assert.Equal(t, 100, building.Health)
			}
		})
	}
}

func TestBuildingsHandler_List(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewBuildingsHandler(store)

	t.Run("list buildings when no city exists", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/city/buildings", nil)
		w := httptest.NewRecorder()

		handler.List(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("list buildings when city exists", func(t *testing.T) {
		// Create a city and add buildings
		store.CreateCity("Test City", "normal")
		store.AddBuilding(models.BuildingFarm)
		store.AddBuilding(models.BuildingHouse)

		req := httptest.NewRequest(http.MethodGet, "/city/buildings", nil)
		w := httptest.NewRecorder()

		handler.List(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var buildings []models.Building
		err := json.Unmarshal(w.Body.Bytes(), &buildings)
		require.NoError(t, err)
		assert.Len(t, buildings, 2)
	})
}

func TestBuildingsHandler_Upgrade(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewBuildingsHandler(store)

	// Create a city and add a building
	store.CreateCity("Test City", "normal")
	building, _ := store.AddBuilding(models.BuildingFarm)

	t.Run("upgrade existing building", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/city/buildings/"+building.ID+"/upgrade", nil)
		w := httptest.NewRecorder()

		// Add chi context with URL parameter
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", building.ID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		handler.Upgrade(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var upgradedBuilding models.Building
		err := json.Unmarshal(w.Body.Bytes(), &upgradedBuilding)
		require.NoError(t, err)
		assert.Equal(t, 2, upgradedBuilding.Level)
	})

	t.Run("upgrade non-existent building", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/city/buildings/nonexistent/upgrade", nil)
		w := httptest.NewRecorder()

		// Add chi context with URL parameter
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "nonexistent")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		handler.Upgrade(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestBuildingsHandler_Repair(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewBuildingsHandler(store)

	// Create a city and add a building
	store.CreateCity("Test City", "normal")
	building, _ := store.AddBuilding(models.BuildingFarm)

	t.Run("repair building with specified amount", func(t *testing.T) {
		request := models.RepairRequest{Amount: 20}
		reqBody, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPatch, "/city/buildings/"+building.ID+"/repair", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Add chi context with URL parameter
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", building.ID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		handler.Repair(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var repairedBuilding models.Building
		err := json.Unmarshal(w.Body.Bytes(), &repairedBuilding)
		require.NoError(t, err)
		assert.Equal(t, 100, repairedBuilding.Health) // Should be capped at 100
	})

	t.Run("repair building with default amount", func(t *testing.T) {
		request := models.RepairRequest{Amount: 0}
		reqBody, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPatch, "/city/buildings/"+building.ID+"/repair", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Add chi context with URL parameter
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", building.ID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		handler.Repair(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestBuildingsHandler_Delete(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewBuildingsHandler(store)

	// Create a city and add a building
	store.CreateCity("Test City", "normal")
	building, _ := store.AddBuilding(models.BuildingFarm)

	t.Run("delete existing building", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/city/buildings/"+building.ID, nil)
		w := httptest.NewRecorder()

		// Add chi context with URL parameter
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", building.ID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		handler.Delete(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify building is deleted
		_, err := store.UpgradeBuilding(building.ID)
		assert.Error(t, err)
	})

	t.Run("delete non-existent building", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/city/buildings/nonexistent", nil)
		w := httptest.NewRecorder()

		// Add chi context with URL parameter
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "nonexistent")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		handler.Delete(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestBuildingsHandler_Effects(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewBuildingsHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/city/buildings/effects", nil)
	w := httptest.NewRecorder()

	handler.Effects(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var effects []models.BuildingEffect
	err := json.Unmarshal(w.Body.Bytes(), &effects)
	require.NoError(t, err)
	assert.Len(t, effects, 4)
}
