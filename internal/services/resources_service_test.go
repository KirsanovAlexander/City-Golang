package services

import (
	"testing"

	"city/internal/models"
	"city/internal/storage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourcesService_Trade(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewResourcesService(store)

	// Create a city first
	store.CreateCity("Test City", "normal")

	tests := []struct {
		name    string
		request models.TradeRequest
		wantErr bool
	}{
		{
			name: "buy food",
			request: models.TradeRequest{
				Resource: "food",
				Amount:   10,
				Price:    2,
			},
			wantErr: false,
		},
		{
			name: "sell energy",
			request: models.TradeRequest{
				Resource: "energy",
				Amount:   -5,
				Price:    3,
			},
			wantErr: false,
		},
		{
			name: "invalid resource",
			request: models.TradeRequest{
				Resource: "invalid",
				Amount:   10,
				Price:    2,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resources, err := service.Trade(tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, resources)
			}
		})
	}
}

func TestResourcesService_GetResourceHistory(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewResourcesService(store)

	history, err := service.GetResourceHistory()
	require.NoError(t, err)
	assert.NotNil(t, history)
}

func TestResourcesService_AdjustResources(t *testing.T) {
	store := storage.NewMemoryStore()
	service := NewResourcesService(store)

	// Create a city first
	store.CreateCity("Test City", "normal")

	tests := []struct {
		name    string
		request models.AdjustResourcesRequest
		wantErr bool
	}{
		{
			name: "adjust food",
			request: models.AdjustResourcesRequest{
				Food: intPtr(200),
			},
			wantErr: false,
		},
		{
			name: "adjust multiple resources",
			request: models.AdjustResourcesRequest{
				Food:   intPtr(150),
				Energy: intPtr(300),
				Money:  intPtr(500),
			},
			wantErr: false,
		},
		{
			name:    "adjust no resources",
			request: models.AdjustResourcesRequest{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resources, err := service.AdjustResources(tt.request)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, resources)
			}
		})
	}
}

func intPtr(i int) *int {
	return &i
}
