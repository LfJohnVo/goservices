package config

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

// FiberConfig func for configuration Fiber app.
func FiberConfig() fiber.Config {
	// Define server settings.
	readTimeoutSecondsCount, _ := strconv.Atoi(GetEnvValue("SERVER_READ_TIMEOUT"))

	// Initialize standard Go html template engine
	engine := html.New("./views", ".html")

	// Return Fiber configuration.
	return fiber.Config{
		Views:       engine,
		ReadTimeout: time.Second * time.Duration(readTimeoutSecondsCount),
	}
}
