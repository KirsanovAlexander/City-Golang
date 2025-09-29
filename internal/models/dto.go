package models

// DTO models - request/response structures

// City DTOs
type CreateCityRequest struct {
	Name       string `json:"name" example:"My City"`
	Difficulty string `json:"difficulty" example:"normal" enums:"easy,normal,hard"`
}

type UpdateSettingsRequest struct {
	Name              *string  `json:"name,omitempty" example:"New City Name"`
	Difficulty        *string  `json:"difficulty,omitempty" example:"hard" enums:"easy,normal,hard"`
	FoodPerFarmer     *int     `json:"foodPerFarmer,omitempty" example:"12"`
	MoneyPerWorker    *int     `json:"moneyPerWorker,omitempty" example:"10"`
	EnergyPerEngineer *int     `json:"energyPerEngineer,omitempty" example:"8"`
	BaseFoodUsePerCap *int     `json:"baseFoodUsePerCap,omitempty" example:"2"`
	HappinessDecay    *float64 `json:"happinessDecay,omitempty" example:"0.7"`
}

// Building DTOs
type BuildRequest struct {
	Type BuildingType `json:"type" example:"farm" enums:"house,farm,factory,powerplant"`
}

type UpgradeRequest struct {
	// Empty for now, but can be extended with upgrade options
}

type RepairRequest struct {
	Amount int `json:"amount" example:"10" minimum:"1" maximum:"100"`
}

type BuildingEffect struct {
	Type   BuildingType `json:"type" example:"farm"`
	Effect string       `json:"effect" example:"+5 food per level per day"`
}

// Citizen DTOs
type CitizenCreateRequest struct {
	Name string `json:"name" example:"John Doe"`
	Job  Job    `json:"job" example:"farmer" enums:"unemployed,farmer,worker,engineer"`
}

type ChangeJobRequest struct {
	Job Job `json:"job" example:"worker" enums:"unemployed,farmer,worker,engineer"`
}

type ChangeHappinessRequest struct {
	Delta float64 `json:"delta" example:"5.0" minimum:"-100" maximum:"100"`
}

type MassAddRequest struct {
	Count  int    `json:"count" example:"5" minimum:"1" maximum:"100"`
	Job    Job    `json:"job" example:"farmer" enums:"unemployed,farmer,worker,engineer"`
	Prefix string `json:"prefix" example:"Citizen"`
}

type JobsStatsResponse struct {
	Unemployed int `json:"unemployed" example:"2"`
	Farmer     int `json:"farmer" example:"5"`
	Worker     int `json:"worker" example:"3"`
	Engineer   int `json:"engineer" example:"1"`
}

// Resource DTOs
type TradeRequest struct {
	Resource string `json:"resource" example:"food" enums:"food,energy,money"`
	Amount   int    `json:"amount" example:"10"` // positive to buy (spend money), negative to sell
	Price    int    `json:"price" example:"2"`   // price per unit in money
}

type AdjustResourcesRequest struct {
	Food   *int `json:"food,omitempty" example:"150"`
	Energy *int `json:"energy,omitempty" example:"200"`
	Money  *int `json:"money,omitempty" example:"300"`
}

// Event DTOs
type CustomEventRequest struct {
	Type      EventType `json:"type" example:"festival" enums:"storm,festival,population_growth,blackout"`
	Message   string    `json:"message" example:"Custom festival event"`
	Delta     Resources `json:"delta"`
	HappDelta float64   `json:"happinessDelta" example:"5.0"`
}

// Error response
type ErrorResponse struct {
	Error string `json:"error" example:"city not found"`
}
