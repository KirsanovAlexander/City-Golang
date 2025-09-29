package services

import (
	"city/internal/models"
	"city/internal/storage"
)

type SimulationService struct {
	store storage.Store
}

func NewSimulationService(store storage.Store) *SimulationService {
	return &SimulationService{store: store}
}

func (s *SimulationService) Tick() (*models.City, error) {
	return s.store.Tick()
}

func (s *SimulationService) RandomEvent() (*models.Event, error) {
	event, err := s.store.RandomEvent()
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (s *SimulationService) CustomEvent(req models.CustomEventRequest) (*models.Event, error) {
	event, err := s.store.CustomEvent(req)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (s *SimulationService) GetEventsHistory() ([]models.Event, error) {
	events := s.store.EventsHistory()
	return events, nil
}

func (s *SimulationService) GetStats() (*models.Stats, error) {
	stats, err := s.store.Stats()
	if err != nil {
		return nil, err
	}
	return &stats, nil
}
