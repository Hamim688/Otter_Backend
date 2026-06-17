package controllers

import (
	"encoding/json"
	"time"

	"Backend/config"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

	oldModeKeamanan := otomatisasi.ModeKeamananAktif

	if err := c.BodyParser(&otomatisasi); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format JSON salah!"})
	}

	config.DB.Save(&otomatisasi)

	// LOGIKA PROFESIONAL: 
	// Jika Mode Keamanan baru saja DIAKTIFKAN, dan saat ini sedang ada pergerakan terdeteksi (TamuGerak = true),
	// langsung nyalakan sirine alarm darurat seketika dan buat notifikasi anomali.
	if otomatisasi.ModeKeamananAktif && !oldModeKeamanan {
		var latestSensor models.SensorLog
		err := config.DB.Order("created_at desc").First(&latestSensor).Error
		if err == nil && latestSensor.TamuGerak {
			var perangkat models.Perangkat
			config.DB.FirstOrCreate(&perangkat, models.Perangkat{ID: 1})
			
			if !perangkat.BuzzerAlrm {
				perangkat.BuzzerAlrm = true
				config.DB.Save(&perangkat)

				// Kirim status terbaru ke ESP32 agar sirine langsung bunyi
				if config.MQTTClient != nil && config.MQTTClient.IsConnected() {
					perangkatJson, _ := json.Marshal(perangkat)
					config.MQTTClient.Publish("otter_smarthome/perangkat", 0, false, perangkatJson)
				}

				pirNotification := models.Notification{
					ID:        uuid.New().String(),
					Title:     "Anomali Terdeteksi",
					Message:   "Ada anomali terdeteksi oleh PIR sensor di Ruang Tamu ketika proteksi baru diaktifkan.",
					Category:  "security",
					Priority:  "critical",
					IsRead:    false,
					Timestamp: time.Now().Format("2006-01-02 15:04:05"),
				}
				config.DB.Create(&pirNotification)
			}
		}
	}

	// Publish ke MQTT biar ESP32 tau rules terbaru
	payload, _ := json.Marshal(otomatisasi)
	config.MQTTClient.Publish("otter_smarthome/otomatisasi", 0, false, payload)

	return c.JSON(fiber.Map{
		"status":  "sukses",
		"message": "Otomatisasi berhasil diupdate!",
		"data":    otomatisasi,
	})
}