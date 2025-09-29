package services

import (
	"testing"

	"city/internal/models"
	"city/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCityService_CreateCity(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewCityService(store)

	tests := []struct {
		name     string
		request  models.CreateCityRequest
		expected string
	}{
		{
			name: "create city with name and difficulty",
			request: models.CreateCityRequest{
				Name:       "Test City",
				Difficulty: "hard",
			},
			expected: "Test City",
		},
		{
			name: "create city with default values",
			request: models.CreateCityRequest{
				Name:       "",
				Difficulty: "",
			},
			expected: "My City",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			city, err := service.CreateCity(tt.request)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, city.Settings.Name)
			assert.NotEmpty(t, city.Settings.Difficulty)
		})
	}
}

func TestCityService_GetCity(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewCityService(store)

	t.Run("get city when no city exists", func(t *testing.T) {
		_, err := service.GetCity()
		assert.Error(t, err)
		assert.Equal(t, ErrNoCity, err)
	})

	t.Run("get city when city exists", func(t *testing.T) {
		store.CreateCity("Test City", "normal")
		city, err := service.GetCity()
		require.NoError(t, err)
		assert.Equal(t, "Test City", city.Settings.Name)
	})
}

func TestCityService_ResetCity(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewCityService(store)

	// Create a city first
	store.CreateCity("Test City", "normal")

	err := service.ResetCity()
	require.NoError(t, err)

	// Verify city is reset
	_, err = store.GetCity()
	assert.Error(t, err)
}

func TestCityService_UpdateSettings(t *testing.T) {
	tests := []struct {
		name         string
		request      models.UpdateSettingsRequest
		expectedName string
		expectedDiff string
	}{
		{
			name: "update city name",
			request: models.UpdateSettingsRequest{
				Name: stringPtr("New City Name"),
			},
			expectedName: "New City Name",
			expectedDiff: "normal",
		},
		{
			name: "update difficulty",
			request: models.UpdateSettingsRequest{
				Difficulty: stringPtr("hard"),
			},
			expectedName: "Test City", // Name should remain unchanged
			expectedDiff: "hard",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fresh store and service for each test
			store := storage.NewMemoryStore()
			service := NewCityService(store)

			// Create a city first
			store.CreateCity("Test City", "normal")

			settings, err := service.UpdateSettings(tt.request)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedName, settings.Name)
			assert.Equal(t, tt.expectedDiff, settings.Difficulty)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
