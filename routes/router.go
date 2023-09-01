package routes

import (
	"goservices/handlers"

	"github.com/gofiber/fiber/v2/middleware/monitor"

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

	api.Get("/timesheet/proyecto", handlers.GetProyectoReport)
	api.Get("/metrics", monitor.New(monitor.Config{Title: "Metrics"}))

	//email send
	api.Post("/emailsend", handlers.PostMail)

}
