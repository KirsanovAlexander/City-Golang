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

func TestCitizensHandler_Create(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewCitizensHandler(store)

	// Create a city first
	store.CreateCity("Test City", "normal")

	tests := []struct {
		name           string
		request        models.CitizenCreateRequest
		expectedStatus int
	}{
		{
			name: "create citizen with name and job",
			request: models.CitizenCreateRequest{
				Name: "John Doe",
				Job:  models.JobFarmer,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "create citizen with default name",
			request: models.CitizenCreateRequest{
				Name: "",
				Job:  models.JobWorker,
			},
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/city/citizens", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Create(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusCreated {
				var citizen models.Citizen
				err := json.Unmarshal(w.Body.Bytes(), &citizen)
				require.NoError(t, err)
				assert.Equal(t, tt.request.Job, citizen.Job)
				assert.NotEmpty(t, citizen.Name)
				assert.Equal(t, 80.0, citizen.Happiness)
			}
		})
	}
}

func TestCitizensHandler_List(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewCitizensHandler(store)

	t.Run("list citizens when no city exists", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/city/citizens", nil)
		w := httptest.NewRecorder()

		handler.List(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("list citizens when city exists", func(t *testing.T) {
		// Create a city and add citizens
		store.CreateCity("Test City", "normal")
		store.AddCitizen("John", models.JobFarmer)
		store.AddCitizen("Jane", models.JobWorker)

		req := httptest.NewRequest(http.MethodGet, "/city/citizens", nil)
		w := httptest.NewRecorder()

		handler.List(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var citizens []models.Citizen
		err := json.Unmarshal(w.Body.Bytes(), &citizens)
		require.NoError(t, err)
		assert.Len(t, citizens, 2)
	})
}

func TestCitizensHandler_ChangeJob(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewCitizensHandler(store)

	// Create a city and add a citizen
	store.CreateCity("Test City", "normal")
	citizen, _ := store.AddCitizen("John", models.JobFarmer)

	t.Run("change job of existing citizen", func(t *testing.T) {
		request := models.ChangeJobRequest{Job: models.JobWorker}
		reqBody, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPatch, "/city/citizens/"+citizen.ID+"/job", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Add chi context with URL parameter
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", citizen.ID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		handler.ChangeJob(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var updatedCitizen models.Citizen
		err := json.Unmarshal(w.Body.Bytes(), &updatedCitizen)
		require.NoError(t, err)
		assert.Equal(t, models.JobWorker, updatedCitizen.Job)
	})

	t.Run("change job of non-existent citizen", func(t *testing.T) {
		request := models.ChangeJobRequest{Job: models.JobWorker}
		reqBody, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPatch, "/city/citizens/nonexistent/job", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Add chi context with URL parameter
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "nonexistent")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		handler.ChangeJob(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestCitizensHandler_ChangeHappiness(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewCitizensHandler(store)

	// Create a city and add a citizen
	store.CreateCity("Test City", "normal")
	citizen, _ := store.AddCitizen("John", models.JobFarmer)

	t.Run("change happiness of existing citizen", func(t *testing.T) {
		request := models.ChangeHappinessRequest{Delta: 10.0}
		reqBody, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPatch, "/city/citizens/"+citizen.ID+"/happiness", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Add chi context with URL parameter
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", citizen.ID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		handler.ChangeHappiness(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var updatedCitizen models.Citizen
		err := json.Unmarshal(w.Body.Bytes(), &updatedCitizen)
		require.NoError(t, err)
		assert.Equal(t, 90.0, updatedCitizen.Happiness)
	})
}

func TestCitizensHandler_Delete(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewCitizensHandler(store)

	// Create a city and add a citizen
	store.CreateCity("Test City", "normal")
	citizen, _ := store.AddCitizen("John", models.JobFarmer)

	t.Run("delete existing citizen", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/city/citizens/"+citizen.ID, nil)
		w := httptest.NewRecorder()

		// Add chi context with URL parameter
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", citizen.ID)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		handler.Delete(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		// Verify citizen is deleted
		_, err := store.ChangeCitizenJob(citizen.ID, models.JobWorker)
		assert.Error(t, err)
	})

	t.Run("delete non-existent citizen", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/city/citizens/nonexistent", nil)
		w := httptest.NewRecorder()

		// Add chi context with URL parameter
		ctx := chi.NewRouteContext()
		ctx.URLParams.Add("id", "nonexistent")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

		handler.Delete(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestCitizensHandler_JobsStats(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewCitizensHandler(store)

	t.Run("get job stats when no city exists", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/city/citizens/jobs", nil)
		w := httptest.NewRecorder()

		handler.JobsStats(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("get job stats when city exists", func(t *testing.T) {
		// Create a city and add citizens with different jobs
		store.CreateCity("Test City", "normal")
		store.AddCitizen("John", models.JobFarmer)
		store.AddCitizen("Jane", models.JobFarmer)
		store.AddCitizen("Bob", models.JobWorker)
		store.AddCitizen("Alice", models.JobEngineer)

		req := httptest.NewRequest(http.MethodGet, "/city/citizens/jobs", nil)
		w := httptest.NewRecorder()

		handler.JobsStats(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var stats models.JobsStatsResponse
		err := json.Unmarshal(w.Body.Bytes(), &stats)
		require.NoError(t, err)
		assert.Equal(t, 2, stats.Farmer)
		assert.Equal(t, 1, stats.Worker)
		assert.Equal(t, 1, stats.Engineer)
		assert.Equal(t, 0, stats.Unemployed)
	})
}

func TestCitizensHandler_MassAdd(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewCitizensHandler(store)

	// Create a city first
	store.CreateCity("Test City", "normal")

	t.Run("mass add citizens", func(t *testing.T) {
		request := models.MassAddRequest{
			Count:  3,
			Job:    models.JobFarmer,
			Prefix: "TestCitizen",
		}
		reqBody, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPost, "/city/citizens/mass-add", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.MassAdd(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var citizens []models.Citizen
		err := json.Unmarshal(w.Body.Bytes(), &citizens)
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
		reqBody, _ := json.Marshal(request)
		req := httptest.NewRequest(http.MethodPost, "/city/citizens/mass-add", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.MassAdd(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var citizens []models.Citizen
		err := json.Unmarshal(w.Body.Bytes(), &citizens)
		require.NoError(t, err)
		assert.Len(t, citizens, 1)
		assert.Equal(t, models.JobWorker, citizens[0].Job)
	})
}
