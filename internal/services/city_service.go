package services

import (
	"errors"

	"city/internal/models"
	"city/internal/storage"
)

var ErrNoCity = errors.New("no city exists")

type CityService struct {
	store storage.Store
}

func NewCityService(store storage.Store) *CityService {
	return &CityService{store: store}
}

func (s *CityService) CreateCity(req models.CreateCityRequest) (*models.City, error) {
	if req.Name == "" {
		req.Name = "My City"
	}
	if req.Difficulty == "" {
		req.Difficulty = "normal"
	}

	city := s.store.CreateCity(req.Name, req.Difficulty)
	return city, nil
}

func (s *CityService) GetCity() (*models.City, error) {
	return s.store.GetCity()
}

func (s *CityService) ResetCity() error {
	s.store.Reset()
	return nil
}

func (s *CityService) UpdateSettings(req models.UpdateSettingsRequest) (*models.CitySettings, error) {
	if err := s.store.UpdateSettings(req); err != nil {
		return nil, err
	}

	city, err := s.store.GetCity()
	if err != nil {
		return nil, err
	}

	return &city.Settings, nil
}
