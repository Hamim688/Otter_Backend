package controllers

import (
	"Backend/config"
	"Backend/models"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

// UpdateSensor memasukkan log sensor baru, memungkinkan simulasi/update sensor dari aplikasi HP.
func UpdateSensor(c *fiber.Ctx) error {
	type UpdateSensorInput struct {
		CahayaAtap      *int     `json:"cahaya_atap"`
		DapurSuhu       *float64 `json:"dapur_suhu"`
		DapurKelembapan *float64 `json:"dapur_kelembapan"`
		DapurFlame      *int     `json:"dapur_flame"`
		KamarSuhu       *float64 `json:"kamar_suhu"`
		KamarKelembapan *float64 `json:"kamar_kelembapan"`
		TamuGerak       *bool    `json:"tamu_gerak"`
	}

	var input UpdateSensorInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format payload tidak valid"})
	}

	// 1. Ambil log sensor terbaru
	var latest models.SensorLog
	config.DB.Order("created_at desc").First(&latest)

	// 2. Salin nilai lama ke log baru
	sensorLog := models.SensorLog{
		CahayaAtap:      latest.CahayaAtap,
		DapurSuhu:       latest.DapurSuhu,
		DapurKelembapan: latest.DapurKelembapan,
		DapurFlame:      latest.DapurFlame,
		KamarSuhu:       latest.KamarSuhu,
		KamarKelembapan: latest.KamarKelembapan,
		TamuGerak:       latest.TamuGerak,
	}

	// 3. Override dengan input baru
	if input.CahayaAtap != nil {
		sensorLog.CahayaAtap = *input.CahayaAtap
	}
	if input.DapurSuhu != nil {
		sensorLog.DapurSuhu = *input.DapurSuhu
	}
	if input.DapurKelembapan != nil {
		sensorLog.DapurKelembapan = *input.DapurKelembapan
	}
	if input.DapurFlame != nil {
		sensorLog.DapurFlame = *input.DapurFlame
	}
	if input.KamarSuhu != nil {
		sensorLog.KamarSuhu = *input.KamarSuhu
	}
	if input.KamarKelembapan != nil {
		sensorLog.KamarKelembapan = *input.KamarKelembapan
	}
	if input.TamuGerak != nil {
		sensorLog.TamuGerak = *input.TamuGerak
	}

	// 4. Simpan ke database
	if err := config.DB.Create(&sensorLog).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal menyimpan data sensor"})
	}

	// 5. Jalankan logika otomatisasi/rules
	var perangkat models.Perangkat
	config.DB.FirstOrCreate(&perangkat, models.Perangkat{ID: 1})

	var otomatisasi models.Otomatisasi
	config.DB.First(&otomatisasi, 1)

	perubahan := false

	// A. Flame Sensor Dapur
	if sensorLog.DapurFlame == 1 && latest.DapurFlame != 1 {
		perangkat.BuzzerAlrm = true
		perangkat.LedMerahDapur = true
		perubahan = true

		// Kirim notifikasi bahaya kebakaran
		fireNotification := models.Notification{
			ID:        uuid.New().String(),
			Title:     "Bahaya Kebakaran!",
			Message:   "Detektor Api Dapur mendeteksi indikasi adanya kebakaran aktif!",
			Category:  "security",
			Priority:  "critical",
			IsRead:    false,
			Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		}
		config.DB.Create(&fireNotification)
	}

	// B. Sensor PIR Ruang Tamu (Hanya jika ModeKeamananAktif = true / rumah kosong)
	if otomatisasi.ModeKeamananAktif && sensorLog.TamuGerak && !latest.TamuGerak {
		perangkat.BuzzerAlrm = true
		perubahan = true

		// Kirim notifikasi anomali terdeteksi
		pirNotification := models.Notification{
			ID:        uuid.New().String(),
			Title:     "Anomali Terdeteksi",
			Message:   "Ada anomali terdeteksi oleh PIR sensor di Ruang Tamu.",
			Category:  "security",
			Priority:  "critical",
			IsRead:    false,
			Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		}
		config.DB.Create(&pirNotification)
	}

	// C. Otomatisasi Suhu Kipas
	if otomatisasi.ModeAutoKipas {
		if sensorLog.KamarSuhu >= otomatisasi.BatasPanasKamar {
			if !perangkat.KipasKamar {
				perangkat.KipasKamar = true
				perubahan = true
			}
			// Tentukan kecepatan kipas berdasarkan selisih suhu
			diff := sensorLog.KamarSuhu - otomatisasi.BatasPanasKamar
			var newSpeed int
			if diff > 4.0 {
				newSpeed = 255
			} else if diff > 2.0 {
				newSpeed = 170
			} else {
				newSpeed = 85
			}
			if perangkat.KecepatanKipas != newSpeed {
				perangkat.KecepatanKipas = newSpeed
				perubahan = true
			}
		} else {
			if perangkat.KipasKamar {
				perangkat.KipasKamar = false
				perangkat.KecepatanKipas = 0
				perubahan = true
			}
		}
	}

	// D. Otomatisasi Cahaya Lampu
	if otomatisasi.ModeAutoLampu {
		isDark := sensorLog.CahayaAtap < otomatisasi.BatasGelapLampu
		if otomatisasi.AutoLampuTamu && perangkat.LampuTamu != isDark {
			perangkat.LampuTamu = isDark
			perubahan = true
		}
		if otomatisasi.AutoLampuKamar && perangkat.LampuKamar != isDark {
			perangkat.LampuKamar = isDark
			perubahan = true
		}
		if otomatisasi.AutoLampuDapur && perangkat.LampuDapur != isDark {
			perangkat.LampuDapur = isDark
			perubahan = true
		}
		if otomatisasi.AutoLampuKamarMandi && perangkat.LampuKamarMandi != isDark {
			perangkat.LampuKamarMandi = isDark
			perubahan = true
		}
	}

	if perubahan {
		config.DB.Save(&perangkat)
		// Kirim update ke MQTT jika client terhubung
		if config.MQTTClient != nil && config.MQTTClient.IsConnected() {
			var payload []byte
			payload, _ = json.Marshal(perangkat)
			config.MQTTClient.Publish("otter_smarthome/perangkat", 0, false, payload)
		}
	}

	return c.JSON(fiber.Map{
		"status": "sukses",
		"sensor": sensorLog,
	})
}
