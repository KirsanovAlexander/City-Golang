package services

import (
	"fmt"

	"city/internal/models"
	"city/internal/storage"
)

type CitizensService struct {
	store storage.Store
}

func NewCitizensService(store storage.Store) *CitizensService {
	return &CitizensService{store: store}
}

func (s *CitizensService) CreateCitizen(req models.CitizenCreateRequest) (*models.Citizen, error) {
	if req.Name == "" {
		req.Name = "Citizen"
	}

	citizen, err := s.store.AddCitizen(req.Name, req.Job)
	if err != nil {
		return nil, err
	}
	return &citizen, nil
}

func (s *CitizensService) ListCitizens() ([]models.Citizen, error) {
	city, err := s.store.GetCity()
	if err != nil {
		return nil, err
	}
	return city.Citizens, nil
}

func (s *CitizensService) ChangeJob(id string, req models.ChangeJobRequest) (*models.Citizen, error) {
	citizen, err := s.store.ChangeCitizenJob(id, req.Job)
	if err != nil {
		return nil, err
	}
	return &citizen, nil
}

func (s *CitizensService) ChangeHappiness(id string, req models.ChangeHappinessRequest) (*models.Citizen, error) {
	citizen, err := s.store.AdjustCitizenHappiness(id, req.Delta)
	if err != nil {
		return nil, err
	}
	return &citizen, nil
}

func (s *CitizensService) DeleteCitizen(id string) error {
	return s.store.RemoveCitizen(id)
}

func (s *CitizensService) GetJobsStats() (*models.JobsStatsResponse, error) {
	city, err := s.store.GetCity()
	if err != nil {
		return nil, err
	}

	stats := &models.JobsStatsResponse{}
	for _, citizen := range city.Citizens {
		switch citizen.Job {
		case models.JobUnemployed:
			stats.Unemployed++
		case models.JobFarmer:
			stats.Farmer++
		case models.JobWorker:
			stats.Worker++
		case models.JobEngineer:
			stats.Engineer++
		}
	}

	return stats, nil
}

func (s *CitizensService) MassAddCitizens(req models.MassAddRequest) ([]models.Citizen, error) {
	if req.Count <= 0 {
		req.Count = 1
	}
	if req.Prefix == "" {
		req.Prefix = "Citizen"
	}

	added := make([]models.Citizen, 0, req.Count)
	for i := 1; i <= req.Count; i++ {
		name := fmt.Sprintf("%s_%d", req.Prefix, i)
		citizen, err := s.store.AddCitizen(name, req.Job)
		if err != nil {
			return nil, err
		}
		added = append(added, citizen)
	}

	return added, nil
}
