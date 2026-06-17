package routes

import (
	"Backend/controllers"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")

	// Tes Koneksi
	api.Get("/ping", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "sukses", "message": "API OTER v4.0 Mengudara!"})
	})

	// === ROUTES PERANGKAT (Sakelar) ===
	api.Get("/perangkat", controllers.GetPerangkat)
	api.Put("/perangkat", controllers.UpdatePerangkat) // Pake PUT karena kita update data

	// === ROUTES OTOMATISASI (Rules) ===
	api.Get("/otomatisasi", controllers.GetOtomatisasi)
	api.Put("/otomatisasi", controllers.UpdateOtomatisasi)
}