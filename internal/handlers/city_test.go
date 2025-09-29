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

func TestCityHandler_Create(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewCityHandler(store)

	tests := []struct {
		name           string
		request        models.CreateCityRequest
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "successful city creation",
			request: models.CreateCityRequest{
				Name:       "Test City",
				Difficulty: "normal",
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "city creation with default values",
			request: models.CreateCityRequest{
				Name:       "",
				Difficulty: "",
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/city", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Create(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if !tt.expectedError {
				var city models.City
				err := json.Unmarshal(w.Body.Bytes(), &city)
				require.NoError(t, err)
				assert.NotEmpty(t, city.Settings.Name)
				assert.NotEmpty(t, city.Settings.Difficulty)
			}
		})
	}
}

func TestCityHandler_Get(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewCityHandler(store)

	t.Run("get city when no city exists", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/city", nil)
		w := httptest.NewRecorder()

		handler.Get(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("get city when city exists", func(t *testing.T) {
		// Create a city first
		store.CreateCity("Test City", "normal")

		req := httptest.NewRequest(http.MethodGet, "/city", nil)
		w := httptest.NewRecorder()

		handler.Get(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var city models.City
		err := json.Unmarshal(w.Body.Bytes(), &city)
		require.NoError(t, err)
		assert.Equal(t, "Test City", city.Settings.Name)
	})
}

func TestCityHandler_Reset(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewCityHandler(store)

	// Create a city first
	store.CreateCity("Test City", "normal")

	req := httptest.NewRequest(http.MethodDelete, "/city/reset", nil)
	w := httptest.NewRecorder()

	handler.Reset(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// Verify city is reset
	_, err := store.GetCity()
	assert.Error(t, err)
}

func TestCityHandler_UpdateSettings(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewCityHandler(store)

	// Create a city first
	store.CreateCity("Test City", "normal")

	tests := []struct {
		name           string
		request        models.UpdateSettingsRequest
		expectedStatus int
	}{
		{
			name: "update city name",
			request: models.UpdateSettingsRequest{
				Name: stringPtr("New City Name"),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "update difficulty",
			request: models.UpdateSettingsRequest{
				Difficulty: stringPtr("hard"),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "update multiple settings",
			request: models.UpdateSettingsRequest{
				Name:       stringPtr("Updated City"),
				Difficulty: stringPtr("easy"),
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPatch, "/city/settings", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.UpdateSettings(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				var settings models.CitySettings
				err := json.Unmarshal(w.Body.Bytes(), &settings)
				require.NoError(t, err)
				if tt.request.Name != nil {
					assert.Equal(t, *tt.request.Name, settings.Name)
				}
				if tt.request.Difficulty != nil {
					assert.Equal(t, *tt.request.Difficulty, settings.Difficulty)
				}
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
