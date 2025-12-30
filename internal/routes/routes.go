package routes

import (
	"github.com/Anezz12/femProject/internal/app"
	"github.com/go-chi/chi"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)
	r.Get("/workouts/{id}", app.WorkoutHandler.HandleGetWorkByID)

	r.Post("/workouts", app.WorkoutHandler.HandleCreateWorkout)

	r.Put("/workouts/{id}", app.WorkoutHandler.HandleUpdateWorkout)

	r.Delete("/workouts/{id}", app.WorkoutHandler.HandleDeleteWorkout)
	r.Post("/users", app.UserHandler.HandleRegisterUser)

	return r
}
