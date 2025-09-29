package models

import "time"

// Domain models - business entities

type BuildingType string

const (
	BuildingHouse      BuildingType = "house"
	BuildingFarm       BuildingType = "farm"
	BuildingFactory    BuildingType = "factory"
	BuildingPowerPlant BuildingType = "powerplant"
)

type Job string

const (
	JobUnemployed Job = "unemployed"
	JobFarmer     Job = "farmer"
	JobWorker     Job = "worker"
	JobEngineer   Job = "engineer"
)

type EventType string

const (
	EventStorm      EventType = "storm"
	EventFestival   EventType = "festival"
	EventPopulation EventType = "population_growth"
	EventBlackout   EventType = "blackout"
)

type CitySettings struct {
	Name              string  `json:"name"`
	Difficulty        string  `json:"difficulty"`
	FoodPerFarmer     int     `json:"foodPerFarmer"`
	MoneyPerWorker    int     `json:"moneyPerWorker"`
	EnergyPerEngineer int     `json:"energyPerEngineer"`
	BaseFoodUsePerCap int     `json:"baseFoodUsePerCap"`
	HappinessDecay    float64 `json:"happinessDecay"`
}

type Resources struct {
	Food   int `json:"food"`
	Energy int `json:"energy"`
	Money  int `json:"money"`
}

type Building struct {
	ID        string       `json:"id"`
	Type      BuildingType `json:"type"`
	Level     int          `json:"level"`
	Health    int          `json:"health"` // 0..100
	CreatedAt time.Time    `json:"createdAt"`
}

type Citizen struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Job       Job       `json:"job"`
	Happiness float64   `json:"happiness"` // 0..100
	CreatedAt time.Time `json:"createdAt"`
}

type City struct {
	Day       int          `json:"day"`
	Settings  CitySettings `json:"settings"`
	Resources Resources    `json:"resources"`
	Buildings []Building   `json:"buildings"`
	Citizens  []Citizen    `json:"citizens"`
}

type Event struct {
	ID        string    `json:"id"`
	Type      EventType `json:"type"`
	Message   string    `json:"message"`
	Delta     Resources `json:"delta"`
	HappDelta float64   `json:"happinessDelta"`
	CreatedAt time.Time `json:"createdAt"`
}

type Stats struct {
	Day           int     `json:"day"`
	Population    int     `json:"population"`
	AvgHappiness  float64 `json:"avgHappiness"`
	FoodBalance   int     `json:"foodBalance"`
	EnergyBalance int     `json:"energyBalance"`
	MoneyBalance  int     `json:"moneyBalance"`
}

type ResourceSnapshot struct {
	Day       int       `json:"day"`
	Resources Resources `json:"resources"`
	Timestamp time.Time `json:"timestamp"`
}
