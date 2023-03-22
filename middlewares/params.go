package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Params(a *fiber.App) {
	a.Use(
		// Add CORS to each route.
		cors.New(cors.Config{
			//cookies passed in frontend
			//this prevent browser blocking
			AllowCredentials: true,
		}),
		// Add simple logger.
		logger.New(),
		//add recover option.
		recover.New(),
	)
}
