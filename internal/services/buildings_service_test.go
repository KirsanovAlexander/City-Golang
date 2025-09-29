package services

import (
	"testing"

	"city/internal/models"
	"city/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildingsService_CreateBuilding(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewBuildingsService(store)

	// Create a city first
	store.CreateCity("Test City", "normal")

	tests := []struct {
		name     string
		request  models.BuildRequest
		expected models.BuildingType
	}{
		{
			name: "create farm building",
			request: models.BuildRequest{
				Type: models.BuildingFarm,
			},
			expected: models.BuildingFarm,
		},
		{
			name: "create house building",
			request: models.BuildRequest{
				Type: models.BuildingHouse,
			},
			expected: models.BuildingHouse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			building, err := service.CreateBuilding(tt.request)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, building.Type)
			assert.Equal(t, 1, building.Level)
			assert.Equal(t, 100, building.Health)
		})
	}
}

func TestBuildingsService_UpgradeBuilding(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewBuildingsService(store)

	// Create a city and add a building
	store.CreateCity("Test City", "normal")
	building, _ := store.AddBuilding(models.BuildingFarm)

	t.Run("upgrade existing building", func(t *testing.T) {
		upgraded, err := service.UpgradeBuilding(building.ID)
		require.NoError(t, err)
		assert.Equal(t, 2, upgraded.Level)
	})

	t.Run("upgrade non-existent building", func(t *testing.T) {
		_, err := service.UpgradeBuilding("nonexistent")
		assert.Error(t, err)
	})
}

func TestBuildingsService_RepairBuilding(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewBuildingsService(store)

	// Create a city and add a building
	store.CreateCity("Test City", "normal")
	building, _ := store.AddBuilding(models.BuildingFarm)

	t.Run("repair building with specified amount", func(t *testing.T) {
		request := models.RepairRequest{Amount: 20}
		repaired, err := service.RepairBuilding(building.ID, request)
		require.NoError(t, err)
		assert.Equal(t, 100, repaired.Health) // Should be capped at 100
	})

	t.Run("repair building with default amount", func(t *testing.T) {
		request := models.RepairRequest{Amount: 0}
		repaired, err := service.RepairBuilding(building.ID, request)
		require.NoError(t, err)
		assert.Equal(t, 100, repaired.Health)
	})
}

func TestBuildingsService_DeleteBuilding(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewBuildingsService(store)

	// Create a city and add a building
	store.CreateCity("Test City", "normal")
	building, _ := store.AddBuilding(models.BuildingFarm)

	t.Run("delete existing building", func(t *testing.T) {
		err := service.DeleteBuilding(building.ID)
		require.NoError(t, err)

		// Verify building is deleted
		_, err = store.UpgradeBuilding(building.ID)
		assert.Error(t, err)
	})

	t.Run("delete non-existent building", func(t *testing.T) {
		err := service.DeleteBuilding("nonexistent")
		assert.Error(t, err)
	})
}

func TestBuildingsService_ListBuildings(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewBuildingsService(store)

	t.Run("list buildings when no city exists", func(t *testing.T) {
		_, err := service.ListBuildings()
		assert.Error(t, err)
	})

	t.Run("list buildings when city exists", func(t *testing.T) {
		// Create a city and add buildings
		store.CreateCity("Test City", "normal")
		store.AddBuilding(models.BuildingFarm)
		store.AddBuilding(models.BuildingHouse)

		buildings, err := service.ListBuildings()
		require.NoError(t, err)
		assert.Len(t, buildings, 2)
	})
}

func TestBuildingsService_GetBuildingEffects(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewBuildingsService(store)

	effects := service.GetBuildingEffects()
	assert.Len(t, effects, 4)

	// Check that all building types are covered
	types := make(map[models.BuildingType]bool)
	for _, effect := range effects {
		types[effect.Type] = true
	}
	assert.True(t, types[models.BuildingFarm])
	assert.True(t, types[models.BuildingFactory])
	assert.True(t, types[models.BuildingPowerPlant])
	assert.True(t, types[models.BuildingHouse])
}
