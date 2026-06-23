package controllers

import (
	"encoding/json"

	"Backend/config"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
)

func GetOtomatisasi(c *fiber.Ctx) error {
	var otomatisasi models.Otomatisasi
	config.DB.FirstOrCreate(&otomatisasi, models.Otomatisasi{ID: 1})
	return c.JSON(otomatisasi)
}

func UpdateOtomatisasi(c *fiber.Ctx) error {
	var otomatisasi models.Otomatisasi

	if err := config.DB.First(&otomatisasi, 1).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Data otomatisasi tidak ditemukan"})
	}

	if err := c.BodyParser(&otomatisasi); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format JSON salah!"})
	}

	config.DB.Save(&otomatisasi)


	// Publish ke MQTT biar ESP32 tau rules terbaru
	payload, _ := json.Marshal(otomatisasi)
	config.MQTTClient.Publish("otter_smarthome/otomatisasi", 0, false, payload)

	// Publish status keamanan ke AI
	statusKeamanan := map[string]bool{"mode_keamanan_aktif": otomatisasi.ModeKeamananAktif}
	statusJson, _ := json.Marshal(statusKeamanan)
	config.MQTTClient.Publish("otter_smarthome/keamanan/mode", 0, true, statusJson)

	return c.JSON(fiber.Map{
		"status":  "sukses",
		"message": "Otomatisasi berhasil diupdate!",
		"data":    otomatisasi,
	})
}