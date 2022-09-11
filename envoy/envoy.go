package main

import (
	"envoy/config"
	"envoy/dbop"
	"envoy/products"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	c, err := config.LoadConfig()

	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	h := dbop.Init(&c)
	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	products.RegisterRoutes(app, h)

	app.Listen(c.Port)
}
