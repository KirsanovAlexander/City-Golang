package router

import (
	"city/internal/handlers"
	"city/internal/storage"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRouter(store storage.Store) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "docs/swagger.json")
	})

	// Initialize handlers
	city := handlers.NewCityHandler(store)
	bld := handlers.NewBuildingsHandler(store)
	cit := handlers.NewCitizensHandler(store)
	res := handlers.NewResourcesHandler(store)
	sim := handlers.NewSimulationHandler(store)

	// City routes
	r.Route("/city", func(r chi.Router) {
		r.Post("/", city.Create)
		r.Get("/", city.Get)
		r.Delete("/reset", city.Reset)
		r.Patch("/settings", city.UpdateSettings)

		// Building routes
		r.Route("/buildings", func(r chi.Router) {
			r.Post("/", bld.Create)
			r.Get("/", bld.List)
			r.Get("/effects", bld.Effects)
			r.Patch("/{id}/upgrade", bld.Upgrade)
			r.Patch("/{id}/repair", bld.Repair)
			r.Delete("/{id}", bld.Delete)
		})

		// Citizen routes
		r.Route("/citizens", func(r chi.Router) {
			r.Post("/", cit.Create)
			r.Get("/", cit.List)
			r.Get("/jobs", cit.JobsStats)
			r.Post("/mass-add", cit.MassAdd)
			r.Patch("/{id}/job", cit.ChangeJob)
			r.Patch("/{id}/happiness", cit.ChangeHappiness)
			r.Delete("/{id}", cit.Delete)
		})

		// Resource routes
		r.Route("/resources", func(r chi.Router) {
			r.Get("/history", res.History)
			r.Patch("/adjust", res.Adjust)
		})

		// Trade route
		r.Post("/trade", res.Trade)

		// Simulation routes
		r.Route("/events", func(r chi.Router) {
			r.Get("/history", sim.EventsHistory)
			r.Post("/random", sim.RandomEvent)
			r.Post("/custom", sim.CustomEvent)
		})

		r.Post("/tick", sim.Tick)
		r.Get("/stats", sim.Stats)
	})

	return r
}
