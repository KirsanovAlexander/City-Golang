package services

import (
	"testing"

	"city/internal/models"
	"city/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimulationService_Tick(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewSimulationService(store)

	t.Run("tick when no city exists", func(t *testing.T) {
		_, err := service.Tick()
		assert.Error(t, err)
	})

	t.Run("tick when city exists", func(t *testing.T) {
		// Create a city with citizens and buildings
		store.CreateCity("Test City", "normal")
		store.AddCitizen("John", models.JobFarmer)
		store.AddBuilding(models.BuildingFarm)

		city, err := service.Tick()
		require.NoError(t, err)
		assert.Equal(t, 1, city.Day)
	})
}

func TestSimulationService_RandomEvent(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewSimulationService(store)

	t.Run("random event when no city exists", func(t *testing.T) {
		_, err := service.RandomEvent()
		assert.Error(t, err)
	})

	t.Run("random event when city exists", func(t *testing.T) {
		// Create a city
		store.CreateCity("Test City", "normal")

		event, err := service.RandomEvent()
		require.NoError(t, err)
		assert.NotEmpty(t, event.ID)
		assert.NotEmpty(t, event.Type)
		assert.NotEmpty(t, event.Message)
	})
}

func TestSimulationService_CustomEvent(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewSimulationService(store)

	// Create a city first
	store.CreateCity("Test City", "normal")

	tests := []struct {
		name    string
		request models.CustomEventRequest
		wantErr bool
	}{
		{
			name: "custom festival event",
			request: models.CustomEventRequest{
				Type:      models.EventFestival,
				Message:   "Custom festival",
				Delta:     models.Resources{Food: -5, Energy: -5, Money: -10},
				HappDelta: 10.0,
			},
			wantErr: false,
		},
		{
			name: "custom storm event",
			request: models.CustomEventRequest{
				Type:      models.EventStorm,
				Message:   "Custom storm",
				Delta:     models.Resources{Food: -10, Energy: -15, Money: -5},
				HappDelta: -5.0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, err := service.CustomEvent(tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.request.Type, event.Type)
				assert.Equal(t, tt.request.Message, event.Message)
				assert.Equal(t, tt.request.HappDelta, event.HappDelta)
			}
		})
	}
}

func TestSimulationService_GetEventsHistory(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewSimulationService(store)

	events, err := service.GetEventsHistory()
	require.NoError(t, err)
	assert.NotNil(t, events)
}

func TestSimulationService_GetStats(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewSimulationService(store)

	t.Run("get stats when no city exists", func(t *testing.T) {
		_, err := service.GetStats()
		assert.Error(t, err)
	})

	t.Run("get stats when city exists", func(t *testing.T) {
		// Create a city with citizens
		store.CreateCity("Test City", "normal")
		store.AddCitizen("John", models.JobFarmer)
		store.AddCitizen("Jane", models.JobWorker)

		stats, err := service.GetStats()
		require.NoError(t, err)
		assert.Equal(t, 0, stats.Day)
		assert.Equal(t, 2, stats.Population)
		assert.Greater(t, stats.AvgHappiness, 0.0)
	})
}
