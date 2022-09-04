package main

import (
	"encoding/json"
	"envoy/dbop"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func test(c *fiber.Ctx) error {
	pl("ok")
	db := dbop.ConnectPostgresClient()
	var Ed []dbop.Lds
	db.Find(&Ed)
	data, _ := json.Marshal(&Ed)
	pl(string(data))
	return nil
}

func main() {
	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	app.Post("/conf/:xds", xdsConfig)
	app.Put("/conf/:xds", xdsConfig)
	app.Delete("/conf/:xds?", xdsConfig)
	app.Post("/v3/discovery:xds", xds)
	app.Get("/", test)
	app.Listen(":8080")
}
