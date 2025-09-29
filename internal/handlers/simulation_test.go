package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"city/internal/models"
	"city/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimulationHandler_Tick(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewSimulationHandler(store)

	t.Run("tick when no city exists", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/city/tick", nil)
		w := httptest.NewRecorder()

		handler.Tick(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("tick when city exists", func(t *testing.T) {
		// Create a city with citizens and buildings
		store.CreateCity("Test City", "normal")
		store.AddCitizen("John", models.JobFarmer)
		store.AddBuilding(models.BuildingFarm)

		req := httptest.NewRequest(http.MethodPost, "/city/tick", nil)
		w := httptest.NewRecorder()

		handler.Tick(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var city models.City
		err := json.Unmarshal(w.Body.Bytes(), &city)
		require.NoError(t, err)
		assert.Equal(t, 1, city.Day)
	})
}

func TestSimulationHandler_RandomEvent(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewSimulationHandler(store)

	t.Run("random event when no city exists", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/city/events/random", nil)
		w := httptest.NewRecorder()

		handler.RandomEvent(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("random event when city exists", func(t *testing.T) {
		// Create a city
		store.CreateCity("Test City", "normal")

		req := httptest.NewRequest(http.MethodPost, "/city/events/random", nil)
		w := httptest.NewRecorder()

		handler.RandomEvent(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var event models.Event
		err := json.Unmarshal(w.Body.Bytes(), &event)
		require.NoError(t, err)
		assert.NotEmpty(t, event.ID)
		assert.NotEmpty(t, event.Type)
		assert.NotEmpty(t, event.Message)
	})
}

func TestSimulationHandler_CustomEvent(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewSimulationHandler(store)

	// Create a city first
	store.CreateCity("Test City", "normal")

	tests := []struct {
		name           string
		request        models.CustomEventRequest
		expectedStatus int
	}{
		{
			name: "custom festival event",
			request: models.CustomEventRequest{
				Type:      models.EventFestival,
				Message:   "Custom festival",
				Delta:     models.Resources{Food: -5, Energy: -5, Money: -10},
				HappDelta: 10.0,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "custom storm event",
			request: models.CustomEventRequest{
				Type:      models.EventStorm,
				Message:   "Custom storm",
				Delta:     models.Resources{Food: -10, Energy: -15, Money: -5},
				HappDelta: -5.0,
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/city/events/custom", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CustomEvent(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				var event models.Event
				err := json.Unmarshal(w.Body.Bytes(), &event)
				require.NoError(t, err)
				assert.Equal(t, tt.request.Type, event.Type)
				assert.Equal(t, tt.request.Message, event.Message)
				assert.Equal(t, tt.request.HappDelta, event.HappDelta)
			}
		})
	}
}

func TestSimulationHandler_EventsHistory(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewSimulationHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/city/events/history", nil)
	w := httptest.NewRecorder()

	handler.EventsHistory(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var events []models.Event
	err := json.Unmarshal(w.Body.Bytes(), &events)
	require.NoError(t, err)
	assert.NotNil(t, events)
}

func TestSimulationHandler_Stats(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewSimulationHandler(store)

	t.Run("get stats when no city exists", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/city/stats", nil)
		w := httptest.NewRecorder()

		handler.Stats(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("get stats when city exists", func(t *testing.T) {
		// Create a city with citizens
		store.CreateCity("Test City", "normal")
		store.AddCitizen("John", models.JobFarmer)
		store.AddCitizen("Jane", models.JobWorker)

		req := httptest.NewRequest(http.MethodGet, "/city/stats", nil)
		w := httptest.NewRecorder()

		handler.Stats(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var stats models.Stats
		err := json.Unmarshal(w.Body.Bytes(), &stats)
		require.NoError(t, err)
		assert.Equal(t, 0, stats.Day)
		assert.Equal(t, 2, stats.Population)
		assert.Greater(t, stats.AvgHappiness, 0.0)
	})
}
