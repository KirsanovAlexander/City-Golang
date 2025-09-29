package storage

import (
	"errors"
	"math/rand/v2"
	"sync"
	"time"

	"city/internal/models"

	"github.com/google/uuid"
)

var ErrNoCity = errors.New("no city exists")

type MemoryStore struct {
	city    *models.City
	events  []models.Event
	resHist []models.ResourceSnapshot
	mu      sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func defaultSettings(name, difficulty string) models.CitySettings {
	s := models.CitySettings{
		Name:              name,
		Difficulty:        difficulty,
		FoodPerFarmer:     10,
		MoneyPerWorker:    8,
		EnergyPerEngineer: 6,
		BaseFoodUsePerCap: 1,
		HappinessDecay:    0.5,
	}
	if difficulty == "hard" {
		s.FoodPerFarmer = 7
		s.MoneyPerWorker = 6
		s.EnergyPerEngineer = 4
		s.BaseFoodUsePerCap = 2
		s.HappinessDecay = 1.0
	}
	return s
}

func (m *MemoryStore) CreateCity(name, difficulty string) *models.City {
	m.mu.Lock()
	defer m.mu.Unlock()
	city := &models.City{
		Day:       0,
		Settings:  defaultSettings(name, difficulty),
		Resources: models.Resources{Food: 100, Energy: 100, Money: 100},
		Buildings: make([]models.Building, 0),
		Citizens:  make([]models.Citizen, 0),
	}
	m.city = city
	m.events = nil
	m.resHist = []models.ResourceSnapshot{{Day: 0, Resources: city.Resources, Timestamp: time.Now()}}
	return city
}

func (m *MemoryStore) GetCity() (*models.City, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.city == nil {
		return nil, ErrNoCity
	}
	cpy := *m.city
	return &cpy, nil
}

func (m *MemoryStore) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.city = nil
	m.events = nil
	m.resHist = nil
}

func (m *MemoryStore) UpdateSettings(req models.UpdateSettingsRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.city == nil {
		return ErrNoCity
	}
	s := &m.city.Settings
	if req.Name != nil {
		s.Name = *req.Name
	}
	if req.Difficulty != nil {
		s.Difficulty = *req.Difficulty
	}
	if req.FoodPerFarmer != nil {
		s.FoodPerFarmer = *req.FoodPerFarmer
	}
	if req.MoneyPerWorker != nil {
		s.MoneyPerWorker = *req.MoneyPerWorker
	}
	if req.EnergyPerEngineer != nil {
		s.EnergyPerEngineer = *req.EnergyPerEngineer
	}
	if req.BaseFoodUsePerCap != nil {
		s.BaseFoodUsePerCap = *req.BaseFoodUsePerCap
	}
	if req.HappinessDecay != nil {
		s.HappinessDecay = *req.HappinessDecay
	}
	return nil
}

// Buildings
func (m *MemoryStore) AddBuilding(t models.BuildingType) (models.Building, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.city == nil {
		return models.Building{}, ErrNoCity
	}
	b := models.Building{ID: uuid.NewString(), Type: t, Level: 1, Health: 100, CreatedAt: time.Now()}
	m.city.Buildings = append(m.city.Buildings, b)
	return b, nil
}

func (m *MemoryStore) UpgradeBuilding(id string) (models.Building, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.city == nil {
		return models.Building{}, ErrNoCity
	}
	for i := range m.city.Buildings {
		if m.city.Buildings[i].ID == id {
			m.city.Buildings[i].Level++
			return m.city.Buildings[i], nil
		}
	}
	return models.Building{}, errors.New("building not found")
}

func (m *MemoryStore) RepairBuilding(id string, amount int) (models.Building, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.city == nil {
		return models.Building{}, ErrNoCity
	}
	for i := range m.city.Buildings {
		if m.city.Buildings[i].ID == id {
			m.city.Buildings[i].Health += amount
			if m.city.Buildings[i].Health > 100 {
				m.city.Buildings[i].Health = 100
			}
			if m.city.Buildings[i].Health < 0 {
				m.city.Buildings[i].Health = 0
			}
			return m.city.Buildings[i], nil
		}
	}
	return models.Building{}, errors.New("building not found")
}

func (m *MemoryStore) RemoveBuilding(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.city == nil {
		return ErrNoCity
	}
	bs := m.city.Buildings
	for i := range bs {
		if bs[i].ID == id {
			m.city.Buildings = append(bs[:i], bs[i+1:]...)
			return nil
		}
	}
	return errors.New("building not found")
}

// Citizens
func (m *MemoryStore) AddCitizen(name string, job models.Job) (models.Citizen, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.city == nil {
		return models.Citizen{}, ErrNoCity
	}
	c := models.Citizen{ID: uuid.NewString(), Name: name, Job: job, Happiness: 80, CreatedAt: time.Now()}
	m.city.Citizens = append(m.city.Citizens, c)
	return c, nil
}

func (m *MemoryStore) ChangeCitizenJob(id string, job models.Job) (models.Citizen, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.city == nil {
		return models.Citizen{}, ErrNoCity
	}
	for i := range m.city.Citizens {
		if m.city.Citizens[i].ID == id {
			m.city.Citizens[i].Job = job
			return m.city.Citizens[i], nil
		}
	}
	return models.Citizen{}, errors.New("citizen not found")
}

func (m *MemoryStore) AdjustCitizenHappiness(id string, delta float64) (models.Citizen, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.city == nil {
		return models.Citizen{}, ErrNoCity
	}
	for i := range m.city.Citizens {
		if m.city.Citizens[i].ID == id {
			m.city.Citizens[i].Happiness += delta
			if m.city.Citizens[i].Happiness > 100 {
				m.city.Citizens[i].Happiness = 100
			}
			if m.city.Citizens[i].Happiness < 0 {
				m.city.Citizens[i].Happiness = 0
			}
			return m.city.Citizens[i], nil
		}
	}
	return models.Citizen{}, errors.New("citizen not found")
}

func (m *MemoryStore) RemoveCitizen(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.city == nil {
		return ErrNoCity
	}
	cs := m.city.Citizens
	for i := range cs {
		if cs[i].ID == id {
			m.city.Citizens = append(cs[:i], cs[i+1:]...)
			return nil
		}
	}
	return errors.New("citizen not found")
}

// Resources
func (m *MemoryStore) Trade(resource string, amount, price int) (models.Resources, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.city == nil {
		return models.Resources{}, ErrNoCity
	}
	cost := amount * price
	switch resource {
	case "food":
		m.city.Resources.Food += amount
	case "energy":
		m.city.Resources.Energy += amount
	case "money":
		m.city.Resources.Money += amount
		cost = 0
	default:
		return models.Resources{}, errors.New("invalid resource")
	}
	m.city.Resources.Money -= cost
	return m.city.Resources, nil
}

func (m *MemoryStore) AdjustResources(req models.AdjustResourcesRequest) (models.Resources, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.city == nil {
		return models.Resources{}, ErrNoCity
	}
	if req.Food != nil {
		m.city.Resources.Food = *req.Food
	}
	if req.Energy != nil {
		m.city.Resources.Energy = *req.Energy
	}
	if req.Money != nil {
		m.city.Resources.Money = *req.Money
	}
	return m.city.Resources, nil
}

func (m *MemoryStore) Snapshot() {
	if m.city == nil {
		return
	}
	m.resHist = append(m.resHist, models.ResourceSnapshot{Day: m.city.Day, Resources: m.city.Resources, Timestamp: time.Now()})
}

func (m *MemoryStore) ResourceHistory() []models.ResourceSnapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]models.ResourceSnapshot, len(m.resHist))
	copy(out, m.resHist)
	return out
}

// Simulation
func (m *MemoryStore) Tick() (*models.City, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.city == nil {
		return nil, ErrNoCity
	}

	s := m.city.Settings
	// Production
	var farmers, workers, engineers int
	for _, c := range m.city.Citizens {
		switch c.Job {
		case models.JobFarmer:
			farmers++
		case models.JobWorker:
			workers++
		case models.JobEngineer:
			engineers++
		}
	}
	m.city.Resources.Food += farmers * s.FoodPerFarmer
	m.city.Resources.Money += workers * s.MoneyPerWorker
	m.city.Resources.Energy += engineers * s.EnergyPerEngineer

	// Buildings passive effects
	for _, b := range m.city.Buildings {
		eff := b.Level
		switch b.Type {
		case models.BuildingFarm:
			m.city.Resources.Food += 5 * eff
		case models.BuildingFactory:
			m.city.Resources.Money += 7 * eff
			m.city.Resources.Energy -= 3 * eff
		case models.BuildingPowerPlant:
			m.city.Resources.Energy += 10 * eff
		case models.BuildingHouse:
			// slight happiness buffer via less decay later
		}
	}

	// Consumption
	pop := len(m.city.Citizens)
	m.city.Resources.Food -= pop * s.BaseFoodUsePerCap
	if m.city.Resources.Food < 0 {
		// hunger penalty
		deficit := -m.city.Resources.Food
		m.city.Resources.Food = 0
		decay := s.HappinessDecay + float64(deficit)/float64(pop+1)
		for i := range m.city.Citizens {
			m.city.Citizens[i].Happiness -= decay
			if m.city.Citizens[i].Happiness < 0 {
				m.city.Citizens[i].Happiness = 0
			}
		}
		// potential population decrease
		if pop > 0 && deficit > pop {
			// remove a random citizen
			idx := rand.IntN(pop)
			m.city.Citizens = append(m.city.Citizens[:idx], m.city.Citizens[idx+1:]...)
		}
	} else {
		// small natural decay
		for i := range m.city.Citizens {
			m.city.Citizens[i].Happiness -= s.HappinessDecay / 2
			if m.city.Citizens[i].Happiness < 0 {
				m.city.Citizens[i].Happiness = 0
			}
		}
	}

	m.city.Day++
	m.resHist = append(m.resHist, models.ResourceSnapshot{Day: m.city.Day, Resources: m.city.Resources, Timestamp: time.Now()})
	cpy := *m.city
	return &cpy, nil
}

func (m *MemoryStore) RandomEvent() (models.Event, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.city == nil {
		return models.Event{}, ErrNoCity
	}
	evTypes := []models.EventType{models.EventStorm, models.EventFestival, models.EventPopulation, models.EventBlackout}
	t := evTypes[rand.IntN(len(evTypes))]
	ev := models.Event{ID: uuid.NewString(), Type: t, CreatedAt: time.Now()}
	switch t {
	case models.EventStorm:
		ev.Message = "A storm hits the city!"
		ev.Delta = models.Resources{Food: -10, Energy: -15, Money: -5}
		ev.HappDelta = -5
		for i := range m.city.Buildings {
			m.city.Buildings[i].Health -= 10
			if m.city.Buildings[i].Health < 0 {
				m.city.Buildings[i].Health = 0
			}
		}
	case models.EventFestival:
		ev.Message = "City festival boosts morale!"
		ev.Delta = models.Resources{Food: -5, Energy: -5, Money: -10}
		ev.HappDelta = 10
	case models.EventPopulation:
		ev.Message = "Population growth!"
		n := 1 + rand.IntN(3)
		for i := 0; i < n; i++ {
			m.city.Citizens = append(m.city.Citizens, models.Citizen{ID: uuid.NewString(), Name: "Newcomer", Job: models.JobUnemployed, Happiness: 70, CreatedAt: time.Now()})
		}
		ev.Delta = models.Resources{}
		ev.HappDelta = 2
	case models.EventBlackout:
		ev.Message = "Power blackout!"
		ev.Delta = models.Resources{Energy: -30}
		ev.HappDelta = -8
	}
	m.city.Resources.Food += ev.Delta.Food
	if m.city.Resources.Food < 0 {
		m.city.Resources.Food = 0
	}
	m.city.Resources.Energy += ev.Delta.Energy
	if m.city.Resources.Energy < 0 {
		m.city.Resources.Energy = 0
	}
	m.city.Resources.Money += ev.Delta.Money
	if m.city.Resources.Money < 0 {
		m.city.Resources.Money = 0
	}
	for i := range m.city.Citizens {
		m.city.Citizens[i].Happiness += ev.HappDelta
		if m.city.Citizens[i].Happiness > 100 {
			m.city.Citizens[i].Happiness = 100
		}
		if m.city.Citizens[i].Happiness < 0 {
			m.city.Citizens[i].Happiness = 0
		}
	}
	m.events = append(m.events, ev)
	return ev, nil
}

func (m *MemoryStore) CustomEvent(req models.CustomEventRequest) (models.Event, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.city == nil {
		return models.Event{}, ErrNoCity
	}
	ev := models.Event{ID: uuid.NewString(), Type: req.Type, Message: req.Message, Delta: req.Delta, HappDelta: req.HappDelta, CreatedAt: time.Now()}
	m.city.Resources.Food += ev.Delta.Food
	if m.city.Resources.Food < 0 {
		m.city.Resources.Food = 0
	}
	m.city.Resources.Energy += ev.Delta.Energy
	if m.city.Resources.Energy < 0 {
		m.city.Resources.Energy = 0
	}
	m.city.Resources.Money += ev.Delta.Money
	if m.city.Resources.Money < 0 {
		m.city.Resources.Money = 0
	}
	for i := range m.city.Citizens {
		m.city.Citizens[i].Happiness += ev.HappDelta
		if m.city.Citizens[i].Happiness > 100 {
			m.city.Citizens[i].Happiness = 100
		}
		if m.city.Citizens[i].Happiness < 0 {
			m.city.Citizens[i].Happiness = 0
		}
	}
	m.events = append(m.events, ev)
	return ev, nil
}

func (m *MemoryStore) EventsHistory() []models.Event {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]models.Event, len(m.events))
	copy(out, m.events)
	return out
}

func (m *MemoryStore) Stats() (models.Stats, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.city == nil {
		return models.Stats{}, ErrNoCity
	}
	var avg float64
	for _, c := range m.city.Citizens {
		avg += c.Happiness
	}
	if len(m.city.Citizens) > 0 {
		avg /= float64(len(m.city.Citizens))
	}
	st := models.Stats{
		Day:           m.city.Day,
		Population:    len(m.city.Citizens),
		AvgHappiness:  avg,
		FoodBalance:   m.city.Resources.Food,
		EnergyBalance: m.city.Resources.Energy,
		MoneyBalance:  m.city.Resources.Money,
	}
	return st, nil
}
