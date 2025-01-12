package main

import (
	"goservices/config"
	"goservices/databases"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-module/carbon/v2"
	"goservices/middlewares"
	"goservices/pkg/utils"
	"goservices/routes"
	"os"
)

func main() {
	carbon.Now().ToDateTimeString()

	// Define Fiber config.
	config := config.FiberConfig()

	//start fiber
	app := fiber.New(config)
	// Middlewares.
	middlewares.Params(app)

	//connect database
	databases.Connect()

	routes.SetupRoutes(app)
	routes.NotFoundRoute(app) // Register route for 404 Error.

	// Start server (with or without graceful shutdown).
	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}

}
