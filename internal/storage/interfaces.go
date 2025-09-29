package storage

import "city/internal/models"

type Store interface {
	// City operations
	CreateCity(name, difficulty string) *models.City
	GetCity() (*models.City, error)
	Reset()
	UpdateSettings(req models.UpdateSettingsRequest) error

	// Building operations
	AddBuilding(t models.BuildingType) (models.Building, error)
	UpgradeBuilding(id string) (models.Building, error)
	RepairBuilding(id string, amount int) (models.Building, error)
	RemoveBuilding(id string) error

	// Citizen operations
	AddCitizen(name string, job models.Job) (models.Citizen, error)
	ChangeCitizenJob(id string, job models.Job) (models.Citizen, error)
	AdjustCitizenHappiness(id string, delta float64) (models.Citizen, error)
	RemoveCitizen(id string) error

	// Resource operations
	Trade(resource string, amount, price int) (models.Resources, error)
	AdjustResources(req models.AdjustResourcesRequest) (models.Resources, error)
	ResourceHistory() []models.ResourceSnapshot

	// Simulation operations
	Tick() (*models.City, error)
	RandomEvent() (models.Event, error)
	CustomEvent(req models.CustomEventRequest) (models.Event, error)
	EventsHistory() []models.Event
	Stats() (models.Stats, error)
}
