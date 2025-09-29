package services

import (
	"city/internal/models"
	"city/internal/storage"
)

type ResourcesService struct {
	store storage.Store
}

func NewResourcesService(store storage.Store) *ResourcesService {
	return &ResourcesService{store: store}
}

func (s *ResourcesService) Trade(req models.TradeRequest) (*models.Resources, error) {
	resources, err := s.store.Trade(req.Resource, req.Amount, req.Price)
	if err != nil {
		return nil, err
	}
	return &resources, nil
}

func (s *ResourcesService) GetResourceHistory() ([]models.ResourceSnapshot, error) {
	history := s.store.ResourceHistory()
	return history, nil
}

func (s *ResourcesService) AdjustResources(req models.AdjustResourcesRequest) (*models.Resources, error) {
	resources, err := s.store.AdjustResources(req)
	if err != nil {
		return nil, err
	}
	return &resources, nil
}
