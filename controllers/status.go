package controllers

import (
	"Backend/config"
	"Backend/models"

	"github.com/gofiber/fiber/v2"
)

// GetStatus mengambil status terbaru seluruh perangkat, otomatisasi, dan sensor
func GetStatus(c *fiber.Ctx) error {
	var perangkat models.Perangkat
	config.DB.FirstOrCreate(&perangkat, models.Perangkat{ID: 1})

	var otomatisasi models.Otomatisasi
	config.DB.FirstOrCreate(&otomatisasi, models.Otomatisasi{ID: 1})

	var sensor models.SensorLog
	// Cari log sensor terbaru, jika tidak ada/error return default sensor kosong
	err := config.DB.Order("created_at desc").First(&sensor).Error

	sensorData := fiber.Map{
		"cahaya_atap":       0,
		"dapur_suhu":        0.0,
		"dapur_kelembapan":  0.0,
		"dapur_flame":       0,
		"kamar_suhu":        0.0,
		"kamar_kelembapan":  0.0,
		"tamu_gerak":        false,
	}

	if err == nil {
		sensorData["cahaya_atap"] = sensor.CahayaAtap
		sensorData["dapur_suhu"] = sensor.DapurSuhu
		sensorData["dapur_kelembapan"] = sensor.DapurKelembapan
		sensorData["dapur_flame"] = sensor.DapurFlame
		sensorData["kamar_suhu"] = sensor.KamarSuhu
		sensorData["kamar_kelembapan"] = sensor.KamarKelembapan
		sensorData["tamu_gerak"] = sensor.TamuGerak
	}

	return c.JSON(fiber.Map{
		"sensor":      sensorData,
		"perangkat":   perangkat,
		"otomatisasi": otomatisasi,
	})
}
