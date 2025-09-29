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

func TestResourcesHandler_Trade(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewResourcesHandler(store)

	// Create a city first
	store.CreateCity("Test City", "normal")

	tests := []struct {
		name           string
		request        models.TradeRequest
		expectedStatus int
	}{
		{
			name: "buy food",
			request: models.TradeRequest{
				Resource: "food",
				Amount:   10,
				Price:    2,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "sell energy",
			request: models.TradeRequest{
				Resource: "energy",
				Amount:   -5,
				Price:    3,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "invalid resource",
			request: models.TradeRequest{
				Resource: "invalid",
				Amount:   10,
				Price:    2,
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/city/trade", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Trade(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				var resources models.Resources
				err := json.Unmarshal(w.Body.Bytes(), &resources)
				require.NoError(t, err)
				assert.NotNil(t, resources)
			}
		})
	}
}

func TestResourcesHandler_History(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewResourcesHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/city/resources/history", nil)
	w := httptest.NewRecorder()

	handler.History(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var history []models.ResourceSnapshot
	err := json.Unmarshal(w.Body.Bytes(), &history)
	require.NoError(t, err)
	assert.NotNil(t, history)
}

func TestResourcesHandler_Adjust(t *testing.T) {
	store := storage.NewMemoryStore()
	handler := NewResourcesHandler(store)

	// Create a city first
	store.CreateCity("Test City", "normal")

	tests := []struct {
		name           string
		request        models.AdjustResourcesRequest
		expectedStatus int
	}{
		{
			name: "adjust food",
			request: models.AdjustResourcesRequest{
				Food: intPtr(200),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "adjust multiple resources",
			request: models.AdjustResourcesRequest{
				Food:   intPtr(150),
				Energy: intPtr(300),
				Money:  intPtr(500),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "adjust no resources",
			request:        models.AdjustResourcesRequest{},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPatch, "/city/resources/adjust", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.Adjust(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				var resources models.Resources
				err := json.Unmarshal(w.Body.Bytes(), &resources)
				require.NoError(t, err)
				assert.NotNil(t, resources)
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}
