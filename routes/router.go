package routes

import (
	"goservices/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	app.Get("/", handlers.Welcome)

	app.Get("/index", func(c *fiber.Ctx) error {
		// Render index template
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	// Middleware
	api := app.Group("/api")

	//Index endpoint
	// api.Post("/register", handlers.Register)
	// api.Post("/login", handlers.Login)
	// api.Get("/authenticated", handlers.AuthenticatedUser)
	// api.Post("/logout", handlers.Logout)
	// api.Post("/forgot", handlers.Forgot)
	// api.Post("/reset", handlers.ResetPassword)
	// api.Get("/metrics", monitor.New(monitor.Config{Title: "Metrics"}))

	//email send
	api.Post("/emailsend", handlers.PostMail)

}
