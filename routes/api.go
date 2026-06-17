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

	// === ROUTES NOTIFIKASI (History Log) ===
	api.Get("/notifications", controllers.GetNotifications)
	api.Put("/notifications/:id/read", controllers.MarkAsRead)
	api.Put("/notifications/read-all", controllers.MarkAllAsRead)
	api.Delete("/notifications/:id", controllers.DeleteNotification)
	api.Delete("/notifications", controllers.ClearAllNotifications)

	// === ROUTES RFID MANAGEMENT ===
	api.Get("/rfid", controllers.GetRfidCards)
	api.Put("/rfid/:uid/approve", controllers.ApproveRfidCard)
	api.Put("/rfid/:uid/status", controllers.UpdateRfidStatus)
	api.Delete("/rfid/:uid", controllers.DeleteRfidCard)

	// === ROUTES HISTORI SENSOR (Grafik) ===
	api.Get("/sensor/history", controllers.GetSensorHistory)
}