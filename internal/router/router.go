package router

import (
	"city/internal/handlers"
	"city/internal/storage"

	"github.com/gin-gonic/gin"
)

func SetupRouter(store *storage.MemoryStore) *gin.Engine {
	r := gin.Default()

	city := handlers.NewCityHandler(store)
	bld := handlers.NewBuildingsHandler(store)
	cit := handlers.NewCitizensHandler(store)
	res := handlers.NewResourcesHandler(store)
	sim := handlers.NewSimulationHandler(store)

	r.POST("/city", city.Create)
	r.GET("/city", city.Get)
	r.DELETE("/city/reset", city.Reset)
	r.PATCH("/city/settings", city.UpdateSettings)

	r.POST("/city/buildings", bld.Create)
	r.PATCH("/city/buildings/:id/upgrade", bld.Upgrade)
	r.DELETE("/city/buildings/:id", bld.Delete)
	r.PATCH("/city/buildings/:id/repair", bld.Repair)
	r.GET("/city/buildings", bld.List)
	r.GET("/city/buildings/effects", bld.Effects)

	r.POST("/city/citizens", cit.Create)
	r.GET("/city/citizens", cit.List)
	r.PATCH("/city/citizens/:id/job", cit.ChangeJob)
	r.PATCH("/city/citizens/:id/happiness", cit.ChangeHappiness)
	r.DELETE("/city/citizens/:id", cit.Delete)
	r.GET("/city/citizens/jobs", cit.JobsStats)
	r.POST("/city/citizens/mass-add", cit.MassAdd)

	r.POST("/city/trade", res.Trade)
	r.GET("/city/resources/history", res.History)
	r.PATCH("/city/resources/adjust", res.Adjust)

	r.POST("/city/tick", sim.Tick)
	r.POST("/city/events/random", sim.RandomEvent)
	r.POST("/city/events/custom", sim.CustomEvent)
	r.GET("/city/events/history", sim.EventsHistory)
	r.GET("/city/stats", sim.Stats)

	return r
}
