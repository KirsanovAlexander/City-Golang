package services

import (
	"city/internal/models"
	"city/internal/storage"
)

type BuildingsService struct {
	store storage.Store
}

func NewBuildingsService(store storage.Store) *BuildingsService {
	return &BuildingsService{store: store}
}

func (s *BuildingsService) CreateBuilding(req models.BuildRequest) (*models.Building, error) {
	building, err := s.store.AddBuilding(req.Type)
	if err != nil {
		return nil, err
	}
	return &building, nil
}

func (s *BuildingsService) UpgradeBuilding(id string) (*models.Building, error) {
	building, err := s.store.UpgradeBuilding(id)
	if err != nil {
		return nil, err
	}
	return &building, nil
}

func (s *BuildingsService) RepairBuilding(id string, req models.RepairRequest) (*models.Building, error) {
	if req.Amount == 0 {
		req.Amount = 10
	}

	building, err := s.store.RepairBuilding(id, req.Amount)
	if err != nil {
		return nil, err
	}
	return &building, nil
}

func (s *BuildingsService) DeleteBuilding(id string) error {
	return s.store.RemoveBuilding(id)
}

func (s *BuildingsService) ListBuildings() ([]models.Building, error) {
	city, err := s.store.GetCity()
	if err != nil {
		return nil, err
	}
	return city.Buildings, nil
}

func (s *BuildingsService) GetBuildingEffects() []models.BuildingEffect {
	return []models.BuildingEffect{
		{Type: models.BuildingFarm, Effect: "+5 food per level per day"},
		{Type: models.BuildingFactory, Effect: "+7 money, -3 energy per level per day"},
		{Type: models.BuildingPowerPlant, Effect: "+10 energy per level per day"},
		{Type: models.BuildingHouse, Effect: "housing/morale, no direct resource"},
	}
}
