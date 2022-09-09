package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	app.Post("/conf/:xds", xdsConfig)
	app.Put("/conf/:xds", xdsConfig)
	app.Delete("/conf/:xds?", xdsConfig)
	app.Post("/v3/discovery:xds", xds)
	app.Listen(":8080")
}
