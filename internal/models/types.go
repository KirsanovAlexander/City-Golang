package models

import (
	"time"
)

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

type EventType string

const (
	EventStorm      EventType = "storm"
	EventFestival   EventType = "festival"
	EventPopulation EventType = "population_growth"
	EventBlackout   EventType = "blackout"
)

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

// Requests
type CreateCityRequest struct {
	Name       string `json:"name"`
	Difficulty string `json:"difficulty"`
}

type UpdateSettingsRequest struct {
	Name              *string  `json:"name"`
	Difficulty        *string  `json:"difficulty"`
	FoodPerFarmer     *int     `json:"foodPerFarmer"`
	MoneyPerWorker    *int     `json:"moneyPerWorker"`
	EnergyPerEngineer *int     `json:"energyPerEngineer"`
	BaseFoodUsePerCap *int     `json:"baseFoodUsePerCap"`
	HappinessDecay    *float64 `json:"happinessDecay"`
}

type BuildRequest struct {
	Type BuildingType `json:"type"` // farm, house, factory, powerplant
}

type UpgradeRequest struct {
}

type RepairRequest struct {
	Amount int `json:"amount"`
}

type CitizenCreateRequest struct {
	Name string `json:"name"`
	Job  Job    `json:"job"`
}

type ChangeJobRequest struct {
	Job Job `json:"job"`
}

type ChangeHappinessRequest struct {
	Delta float64 `json:"delta"`
}

type MassAddRequest struct {
	Count  int    `json:"count"`
	Job    Job    `json:"job"`
	Prefix string `json:"prefix"`
}

type TradeRequest struct {
	Resource string `json:"resource"` // food, energy, money
	Amount   int    `json:"amount"`   // positive to buy (spend money), negative to sell
	Price    int    `json:"price"`    // price per unit in money
}

type AdjustResourcesRequest struct {
	Food   *int `json:"food"`
	Energy *int `json:"energy"`
	Money  *int `json:"money"`
}

type CustomEventRequest struct {
	Type      EventType `json:"type"`
	Message   string    `json:"message"`
	Delta     Resources `json:"delta"`
	HappDelta float64   `json:"happinessDelta"`
}

type ResourceSnapshot struct {
	Day       int       `json:"day"`
	Resources Resources `json:"resources"`
	Timestamp time.Time `json:"timestamp"`
}
