package routes

import (
	"fiber-boilerplate/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
)

func SetupRoutes(app *fiber.App) {

	app.Get("/", handlers.Welcome)
	// Middleware
	api := app.Group("/api")

	//Index endpoint
	api.Post("/register", handlers.Register)
	api.Post("/login", handlers.Login)
	api.Get("/authenticated", handlers.AuthenticatedUser)
	api.Post("/logout", handlers.Logout)
	api.Post("/forgot", handlers.Forgot)
	api.Post("/reset", handlers.ResetPassword)
	api.Get("/metrics", monitor.New(monitor.Config{Title: "Metrics"}))

	//Dogs
	api.Get("/dogs", handlers.GetDogs)
	api.Get("/dogs/:id", handlers.GetDog)
	api.Post("/dogs", handlers.AddDog)
	api.Put("/dogs/:id", handlers.UpdateDog)
	api.Delete("/dogs/:id", handlers.RemoveDog)

}
