package controllers

import (
	"Backend/config"
	"time"

	"github.com/gofiber/fiber/v2"
)

type HistoryResponse struct {
	Timestamp       string  `json:"timestamp"`
	DapurSuhu       float64 `json:"dapur_suhu"`
	DapurKelembapan float64 `json:"dapur_kelembapan"`
	KamarSuhu       float64 `json:"kamar_suhu"`
	KamarKelembapan float64 `json:"kamar_kelembapan"`
}

// GetSensorHistory mengambil data rata-rata sensor suhu & kelembapan terkelompok untuk grafik di HP
func GetSensorHistory(c *fiber.Ctx) error {
	timeRange := c.Query("range", "24h")
	var results []HistoryResponse
	var query string

	now := time.Now()
	var startTime time.Time

	if timeRange == "7d" {
		startTime = now.AddDate(0, 0, -7)
		// Group by day di PostgreSQL (YYYY-MM-DD)
		query = `SELECT TO_CHAR(created_at, 'YYYY-MM-DD') as timestamp, 
		                ROUND(AVG(dapur_suhu)::numeric, 1) as dapur_suhu,
		                ROUND(AVG(dapur_kelembapan)::numeric, 1) as dapur_kelembapan,
		                ROUND(AVG(kamar_suhu)::numeric, 1) as kamar_suhu,
		                ROUND(AVG(kamar_kelembapan)::numeric, 1) as kamar_kelembapan
		         FROM sensor_logs 
		         WHERE created_at >= ? 
		         GROUP BY TO_CHAR(created_at, 'YYYY-MM-DD')
		         ORDER BY timestamp ASC`
	} else if timeRange == "30d" {
		startTime = now.AddDate(0, 0, -30)
		// Group by day di PostgreSQL (YYYY-MM-DD)
		query = `SELECT TO_CHAR(created_at, 'YYYY-MM-DD') as timestamp, 
		                ROUND(AVG(dapur_suhu)::numeric, 1) as dapur_suhu,
		                ROUND(AVG(dapur_kelembapan)::numeric, 1) as dapur_kelembapan,
		                ROUND(AVG(kamar_suhu)::numeric, 1) as kamar_suhu,
		                ROUND(AVG(kamar_kelembapan)::numeric, 1) as kamar_kelembapan
		         FROM sensor_logs 
		         WHERE created_at >= ? 
		         GROUP BY TO_CHAR(created_at, 'YYYY-MM-DD')
		         ORDER BY timestamp ASC`
	} else { // default 24h
		startTime = now.Add(-24 * time.Hour)
		// Group by hour di PostgreSQL (YYYY-MM-DD HH24:00)
		query = `SELECT TO_CHAR(created_at, 'YYYY-MM-DD HH24:00') as timestamp, 
		                ROUND(AVG(dapur_suhu)::numeric, 1) as dapur_suhu,
		                ROUND(AVG(dapur_kelembapan)::numeric, 1) as dapur_kelembapan,
		                ROUND(AVG(kamar_suhu)::numeric, 1) as kamar_suhu,
		                ROUND(AVG(kamar_kelembapan)::numeric, 1) as kamar_kelembapan
		         FROM sensor_logs 
		         WHERE created_at >= ? 
		         GROUP BY TO_CHAR(created_at, 'YYYY-MM-DD HH24:00')
		         ORDER BY timestamp ASC`
	}

	if err := config.DB.Raw(query, startTime).Scan(&results).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal memproses data histori sensor"})
	}

	// Jika data kosong, kembalikan array kosong agar aplikasi Flutter tidak crash
	if results == nil {
		results = []HistoryResponse{}
	}

	return c.JSON(results)
}
