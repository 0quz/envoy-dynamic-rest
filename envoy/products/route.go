package products

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	h := &handler{
		DB: db,
	}

	app.Post("/conf/:xds", h.XdsConfig)
	app.Put("/conf/:xds", h.XdsConfig)
	app.Delete("/conf/:xds", h.XdsConfig)
	app.Post("/v3/discovery:xds", h.Xds)
}
