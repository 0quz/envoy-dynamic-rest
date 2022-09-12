package products

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// DB handler
type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(app *fiber.App, db *gorm.DB) {
	h := &handler{
		DB: db,
	}
	// Endpoints definition and send DB handler
	app.Post("/conf/:xds", h.XdsConfig)
	app.Put("/conf/:xds", h.XdsConfig)
	app.Delete("/conf/:xds", h.XdsConfig)
	app.Post("/v3/discovery:xds", h.Xds)
}
