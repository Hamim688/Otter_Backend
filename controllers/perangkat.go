package controllers

import (
	"encoding/json"
	"Backend/config"
	"Backend/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Ambil Status Perangkat Saat Ini
func GetPerangkat(c *fiber.Ctx) error {
	var perangkat models.Perangkat
	
	// Cari data dengan ID 1. Kalau belum ada, bikin baru (FirstOrCreate)
	config.DB.FirstOrCreate(&perangkat, models.Perangkat{ID: 1})

	return c.JSON(perangkat)
}

// Update Status Perangkat dari Flutter
func UpdatePerangkat(c *fiber.Ctx) error {
	var perangkat models.Perangkat

	// 1. Ambil data lama sebelum di-update untuk pembanding
	var oldPerangkat models.Perangkat
	if err := config.DB.First(&oldPerangkat, 1).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Data perangkat tidak ditemukan"})
	}

	// 2. Timpa data lama dengan data baru dari JSON Flutter
	perangkat = oldPerangkat // Jaga agar ID tetap 1
	if err := c.BodyParser(&perangkat); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format JSON salah!"})
	}



	// 4. Save ke PostgreSQL
	config.DB.Save(&perangkat)

	// 5. PUBLISH KE MQTT BIAR ESP32 LANGSUNG GERAK!
	payload, _ := json.Marshal(perangkat)
	config.MQTTClient.Publish("otter_smarthome/perangkat", 0, false, payload)

	return c.JSON(fiber.Map{
		"status":  "sukses",
		"message": "Perangkat berhasil diupdate dan perintah MQTT dikirim!",
		"data":    perangkat,
	})
}