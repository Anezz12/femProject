package routes

import (
	"github.com/Anezz12/femProject/internal/app"
	"github.com/go-chi/chi"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(app.Middleware.Authenticate)

		r.Get("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleGetWorkByID))

		r.Post("/workouts", app.Middleware.RequireUser(app.WorkoutHandler.HandleCreateWorkout))

		r.Put("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleUpdateWorkout))

		r.Delete("/workouts/{id}", app.Middleware.RequireUser(app.WorkoutHandler.HandleDeleteWorkout))
	})

	r.Get("/health", app.HealthCheck)
	r.Post("/users", app.UserHandler.HandleRegisterUser)
	r.Post("/tokens/authentication", app.TokenHandler.HandlerCreateToken)

	return r
}
