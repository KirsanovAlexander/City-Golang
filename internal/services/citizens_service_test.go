package services

import (
	"testing"

	"city/internal/models"
	"city/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCitizensService_CreateCitizen(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewCitizensService(store)

	// Create a city first
	store.CreateCity("Test City", "normal")

	tests := []struct {
		name     string
		request  models.CitizenCreateRequest
		expected string
	}{
		{
			name: "create citizen with name and job",
			request: models.CitizenCreateRequest{
				Name: "John Doe",
				Job:  models.JobFarmer,
			},
			expected: "John Doe",
		},
		{
			name: "create citizen with default name",
			request: models.CitizenCreateRequest{
				Name: "",
				Job:  models.JobWorker,
			},
			expected: "Citizen",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			citizen, err := service.CreateCitizen(tt.request)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, citizen.Name)
			assert.Equal(t, tt.request.Job, citizen.Job)
			assert.Equal(t, 80.0, citizen.Happiness)
		})
	}
}

func TestCitizensService_ListCitizens(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewCitizensService(store)

	t.Run("list citizens when no city exists", func(t *testing.T) {
		_, err := service.ListCitizens()
		assert.Error(t, err)
	})

	t.Run("list citizens when city exists", func(t *testing.T) {
		// Create a city and add citizens
		store.CreateCity("Test City", "normal")
		store.AddCitizen("John", models.JobFarmer)
		store.AddCitizen("Jane", models.JobWorker)

		citizens, err := service.ListCitizens()
		require.NoError(t, err)
		assert.Len(t, citizens, 2)
	})
}

func TestCitizensService_ChangeJob(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewCitizensService(store)

	// Create a city and add a citizen
	store.CreateCity("Test City", "normal")
	citizen, _ := store.AddCitizen("John", models.JobFarmer)

	t.Run("change job of existing citizen", func(t *testing.T) {
		request := models.ChangeJobRequest{Job: models.JobWorker}
		updated, err := service.ChangeJob(citizen.ID, request)
		require.NoError(t, err)
		assert.Equal(t, models.JobWorker, updated.Job)
	})

	t.Run("change job of non-existent citizen", func(t *testing.T) {
		request := models.ChangeJobRequest{Job: models.JobWorker}
		_, err := service.ChangeJob("nonexistent", request)
		assert.Error(t, err)
	})
}

func TestCitizensService_ChangeHappiness(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewCitizensService(store)

	// Create a city and add a citizen
	store.CreateCity("Test City", "normal")
	citizen, _ := store.AddCitizen("John", models.JobFarmer)

	t.Run("change happiness of existing citizen", func(t *testing.T) {
		request := models.ChangeHappinessRequest{Delta: 10.0}
		updated, err := service.ChangeHappiness(citizen.ID, request)
		require.NoError(t, err)
		assert.Equal(t, 90.0, updated.Happiness)
	})
}

func TestCitizensService_DeleteCitizen(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewCitizensService(store)

	// Create a city and add a citizen
	store.CreateCity("Test City", "normal")
	citizen, _ := store.AddCitizen("John", models.JobFarmer)

	t.Run("delete existing citizen", func(t *testing.T) {
		err := service.DeleteCitizen(citizen.ID)
		require.NoError(t, err)

		// Verify citizen is deleted
		_, err = store.ChangeCitizenJob(citizen.ID, models.JobWorker)
		assert.Error(t, err)
	})

	t.Run("delete non-existent citizen", func(t *testing.T) {
		err := service.DeleteCitizen("nonexistent")
		assert.Error(t, err)
	})
}

func TestCitizensService_GetJobsStats(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewCitizensService(store)

	t.Run("get job stats when no city exists", func(t *testing.T) {
		_, err := service.GetJobsStats()
		assert.Error(t, err)
	})

	t.Run("get job stats when city exists", func(t *testing.T) {
		// Create a city and add citizens with different jobs
		store.CreateCity("Test City", "normal")
		store.AddCitizen("John", models.JobFarmer)
		store.AddCitizen("Jane", models.JobFarmer)
		store.AddCitizen("Bob", models.JobWorker)
		store.AddCitizen("Alice", models.JobEngineer)

		stats, err := service.GetJobsStats()
		require.NoError(t, err)
		assert.Equal(t, 2, stats.Farmer)
		assert.Equal(t, 1, stats.Worker)
		assert.Equal(t, 1, stats.Engineer)
		assert.Equal(t, 0, stats.Unemployed)
	})
}

func TestCitizensService_MassAddCitizens(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewCitizensService(store)

	// Create a city first
	store.CreateCity("Test City", "normal")

	t.Run("mass add citizens", func(t *testing.T) {
		request := models.MassAddRequest{
			Count:  3,
			Job:    models.JobFarmer,
			Prefix: "TestCitizen",
		}
		citizens, err := service.MassAddCitizens(request)
		require.NoError(t, err)
		assert.Len(t, citizens, 3)
		for _, citizen := range citizens {
			assert.Equal(t, models.JobFarmer, citizen.Job)
		}
	})

	t.Run("mass add with default values", func(t *testing.T) {
		request := models.MassAddRequest{
			Count:  0, // Should default to 1
			Job:    models.JobWorker,
			Prefix: "", // Should default to "Citizen"
		}
		citizens, err := service.MassAddCitizens(request)
		require.NoError(t, err)
		assert.Len(t, citizens, 1)
		assert.Equal(t, models.JobWorker, citizens[0].Job)
	})
}
