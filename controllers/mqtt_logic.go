package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"Backend/config"
	"Backend/models"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

type SensorPayload struct {
	CahayaAtap      int     `json:"cahaya_atap"`
	DapurSuhu       float64 `json:"dapur_suhu"`
	DapurKelembapan float64 `json:"dapur_kelembapan"`
	DapurFlame      int     `json:"dapur_flame"`
	KamarSuhu       float64 `json:"kamar_suhu"`
	KamarKelembapan float64 `json:"kamar_kelembapan"`
	TamuGerak       bool    `json:"tamu_gerak"`
}

type RfidPayload struct {
	UID string `json:"uid"`
}

var MessagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	payloadStr := string(msg.Payload())

	fmt.Printf("[MQTT MASUK] Topik: %s \n[PAYLOAD]: %s\n\n", topic, payloadStr)

	// ========================================================
	// 1. LOGIKA UTK DATA SENSOR (otter_smarthome/sensor)
	// ========================================================
	if topic == "otter_smarthome/sensor" {
		var payload SensorPayload
		if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
			fmt.Println("[MQTT ERROR] Gagal parse JSON sensor:", err)
			return
		}

		// Konversi CahayaAtap (LDR) ke Persentase (0-100%)
		// Asumsi ESP32 ADC Max = 4095. Semakin terang, nilai persentase semakin tinggi (mendekati 100%).
		cahayaPercentage := 100 - int((float64(payload.CahayaAtap)/4095.0)*100)
		if cahayaPercentage < 0 {
			cahayaPercentage = 0
		} else if cahayaPercentage > 100 {
			cahayaPercentage = 100
		}

		// A. Simpan log sensor ke PostgreSQL
		sensorLog := models.SensorLog{
			CahayaAtap:      cahayaPercentage,
			DapurSuhu:       payload.DapurSuhu,
			DapurKelembapan: payload.DapurKelembapan,
			DapurFlame:      payload.DapurFlame,
			KamarSuhu:       payload.KamarSuhu,
			KamarKelembapan: payload.KamarKelembapan,
			TamuGerak:       payload.TamuGerak,
		}
		if err := config.DB.Create(&sensorLog).Error; err != nil {
			fmt.Println("[DATABASE ERROR] Gagal menyimpan log sensor:", err)
		}

		// C. Logika Otomatisasi (Suhu & Cahaya)
		var aturan models.Otomatisasi
		if err := config.DB.First(&aturan, 1).Error; err == nil {
			perubahan := false
			var perangkat models.Perangkat
			config.DB.FirstOrCreate(&perangkat, models.Perangkat{ID: 1})

      // 1. Cek Kipas Auto
      if aturan.ModeAutoKipas && payload.KamarSuhu > aturan.BatasPanasKamar {
        if !perangkat.KipasKamar {
          perangkat.KipasKamar = true
          perangkat.KecepatanKipas = 255 // Kecepatan penuh saat auto nyala
          perubahan = true
          fmt.Println("[AUTO] Kamar Kepanasan! Kipas otomatis MENYALA.")
        }
      }

			// 2. Cek Lampu Auto
			if aturan.ModeAutoLampu && cahayaPercentage < aturan.BatasGelapLampu {
				if !perangkat.LampuKamar {
					perangkat.LampuKamar = true // Lampu nyala!
					perubahan = true
					fmt.Println("[AUTO] Kamar Gelap! Lampu kamar otomatis MENYALA.")
				}
			}

      // 3. Kalau sistem ngerubah kipas/lampu, simpan ke DB & lapor ke ESP32
      if perubahan {
        config.DB.Save(&perangkat)
        perangkatJson, _ := json.Marshal(perangkat)
        config.MQTTClient.Publish("otter_smarthome/perangkat", 0, false, perangkatJson)
      }
    }
	}

	// ========================================================
	// 2. LOGIKA RFID SCAN (otter_smarthome/rfid_terdaftar/scan_terbaru)
	// ========================================================
	if topic == "otter_smarthome/rfid_terdaftar/scan_terbaru" {
		var payload RfidPayload
		if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
			// Jika payload berupa UID string polos, pasangkan manual
			payload.UID = payloadStr
		}

		if payload.UID == "" {
			fmt.Println("[MQTT ERROR] UID RFID Kosong.")
			return
		}

		var card models.RfidCard
		result := config.DB.Where("uid = ?", payload.UID).First(&card)

		if result.Error == nil {
			// A. KARTU TERDAFTAR DAN AKTIF
			if card.Status == "aktif" {
				fmt.Printf("[RFID] Akses Diterima untuk: %s (%s)\n", card.NamaPemilik, card.UID)

				// Unlock pintu di database
				var perangkat models.Perangkat
				config.DB.FirstOrCreate(&perangkat, models.Perangkat{ID: 1})
				perangkat.KunciPintuRfid = false // False = Pintu terbuka
				config.DB.Save(&perangkat)

				// Set ModeKeamananAktif = false (karena ada orang masuk)
				var otomatisasi models.Otomatisasi
				if err := config.DB.First(&otomatisasi, 1).Error; err == nil {
					otomatisasi.ModeKeamananAktif = false
					config.DB.Save(&otomatisasi)

					// Kode Untuk AI 
					statusKeamanan := map[string]bool{"mode_keamanan_aktif": false}
          statusJson, _ := json.Marshal(statusKeamanan)
          config.MQTTClient.Publish("otter_smarthome/keamanan/mode", 0, true, statusJson)
				}

				// Kirim status pintu terbaru ke ESP32 via MQTT agar Servo berputar
				perangkatJson, _ := json.Marshal(perangkat)
				config.MQTTClient.Publish("otter_smarthome/perangkat", 0, false, perangkatJson)

				// [FITUR BARU] Auto-Lock Pintu Setelah 5 Detik
				time.AfterFunc(5*time.Second, func() {
					fmt.Println("[AUTO-LOCK] 5 detik berlalu. Mengunci pintu kembali...")
					
					var p models.Perangkat
					if err := config.DB.First(&p, 1).Error; err == nil {
						if !p.KunciPintuRfid { // Jika masih terbuka
							p.KunciPintuRfid = true // Kunci kembali
							config.DB.Save(&p)
							
							pJson, _ := json.Marshal(p)
							config.MQTTClient.Publish("otter_smarthome/perangkat", 0, false, pJson)
							fmt.Println("[AUTO-LOCK] Pintu berhasil dikunci.")
						}
					}
				})

				// Kirim feedback respons sukses ke topik khusus RFID
				responsePayload := map[string]string{
					"uid":          card.UID,
					"status":       "aktif",
					"nama_pemilik": card.NamaPemilik,
				}
				respJson, _ := json.Marshal(responsePayload)
				config.MQTTClient.Publish("otter_smarthome/rfid/response", 0, false, respJson)

				// Buat notifikasi sukses masuk
				newNotification := models.Notification{
					ID:        uuid.New().String(),
					Title:     "Akses Pintu Utama",
					Message:   fmt.Sprintf("Akses masuk berhasil dibuka oleh %s via RFID.", card.NamaPemilik),
					Category:  "security",
					Priority:  "info",
					IsRead:    false,
					Timestamp: time.Now().Format("2006-01-02 15:04:05"),
				}
				config.DB.Create(&newNotification)

			} else if card.Status == "menunggu" {
				// B. KARTU SEDANG MENUNGGU PERSETUJUAN
				fmt.Printf("[RFID] Akses Tertunda: Kartu %s masih menunggu persetujuan.\n", card.UID)
				
				responsePayload := map[string]string{
					"uid":    card.UID,
					"status": "menunggu",
				}
				respJson, _ := json.Marshal(responsePayload)
				config.MQTTClient.Publish("otter_smarthome/rfid/response", 0, false, respJson)

			} else {
				// C. KARTU DINONAKTIFKAN
				fmt.Printf("[RFID] Akses Ditolak: Kartu %s dinonaktifkan.\n", card.UID)
				
				responsePayload := map[string]string{
					"uid":    card.UID,
					"status": "nonaktif",
				}
				respJson, _ := json.Marshal(responsePayload)
				config.MQTTClient.Publish("otter_smarthome/rfid/response", 0, false, respJson)
			}
		} else {
			// D. KARTU BARU/ASING (BELUM TERDAFTAR DI DATABASE)
			fmt.Printf("[RFID] Kartu Baru Terdeteksi: %s. Menyimpan ke database sebagai 'menunggu'...\n", payload.UID)

			// Simpan kartu baru dengan status 'menunggu' agar bisa disetujui dari HP
			newCard := models.RfidCard{
				UID:         payload.UID,
				NamaPemilik: "Unknown Card",
				Status:      "menunggu",
			}
			config.DB.Create(&newCard)

			// Kirim response status menunggu ke ESP32 agar berbunyi bip alarm penolakan
			responsePayload := map[string]string{
				"uid":    payload.UID,
				"status": "menunggu",
			}
			respJson, _ := json.Marshal(responsePayload)
			config.MQTTClient.Publish("otter_smarthome/rfid/response", 0, false, respJson)

			// Buat notifikasi peringatan kartu asing terdeteksi
			warnNotification := models.Notification{
				ID:        uuid.New().String(),
				Title:     "Peringatan RFID Asing",
				Message:   fmt.Sprintf("Kartu RFID asing dengan UID %s terdeteksi menempel pada alat.", payload.UID),
				Category:  "security",
				Priority:  "warning",
				IsRead:    false,
				Timestamp: time.Now().Format("2006-01-02 15:04:05"),
			}
			config.DB.Create(&warnNotification)
		}
	}

	// ========================================================
	// 3. LOGIKA AI ALERT (otter_smarthome/ai_alert)
	// ========================================================
	if topic == "otter_smarthome/ai_alert" {
		fmt.Printf("[🚨 ALARM AI] %s\n", payloadStr)

		var alertData map[string]string
		if err := json.Unmarshal(msg.Payload(), &alertData); err == nil {

			// A. Simpan peringatan ke database biar HP Rafa dapet notif
			newNotification := models.Notification{
				ID:        uuid.New().String(),
				Title:     alertData["status"], // Isinya "BAHAYA" dari Python
				Message:   alertData["pesan"],  // Isinya "AI mendeteksi anomali..."
				Category:  "security",
				Priority:  "critical",
				IsRead:    false,
				Timestamp: alertData["timestamp"],
			}
			config.DB.Create(&newNotification)

			// B. Otomatis nyalain Sirine (Buzzer) di ESP32
			var perangkat models.Perangkat
			if err := config.DB.FirstOrCreate(&perangkat, models.Perangkat{ID: 1}).Error; err == nil {
				if !perangkat.BuzzerAlrm {
					perangkat.BuzzerAlrm = true
					config.DB.Save(&perangkat)

					// Tembak perintah ke ESP32 buat bunyiin buzzer
					perangkatJson, _ := json.Marshal(perangkat)
					config.MQTTClient.Publish("otter_smarthome/perangkat", 0, false, perangkatJson)
					
					// [FITUR BARU] Kedip LED Merah Dapur (1-0-1-0) setiap 1 detik saat alarm aktif
					go func() {
						ledState := true
						for {
							// Cek apakah alarm masih aktif di database
							var p models.Perangkat
							if err := config.DB.First(&p, 1).Error; err != nil || !p.BuzzerAlrm {
								// Jika alarm dimatikan (Disarm), pastikan LED Merah mati lalu keluar dari loop
								p.LedMerahDapur = false
								config.DB.Save(&p)
								pJson, _ := json.Marshal(p)
								config.MQTTClient.Publish("otter_smarthome/perangkat", 0, false, pJson)
								break
							}

							// Toggle state LED
							p.LedMerahDapur = ledState
							config.DB.Save(&p)
							
							pJson, _ := json.Marshal(p)
							config.MQTTClient.Publish("otter_smarthome/perangkat", 0, false, pJson)

							ledState = !ledState // balikkan status untuk iterasi berikutnya (kedip)
							time.Sleep(1 * time.Second)
						}
					}()
				}
			}
		}
	}
}