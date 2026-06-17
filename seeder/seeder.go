package seeder

import (
	"log"
	"Backend/models"

	"gorm.io/gorm"
)

// Seed mempopulasikan database dengan data awal jika tabel kosong
func Seed(db *gorm.DB) {
	// === SEED DATA DUMMY AKUN ===
	var userCount int64
	db.Model(&models.User{}).Count(&userCount)
	if userCount == 0 {
		dummyUser := models.User{
			Username:    "admin",
			Password:    "1234",
			DisplayName: "Mimah Dudim",
			Role:        "Administrator Rumah Pintar",
			AvatarURL:   "",
		}
		if err := db.Create(&dummyUser).Error; err == nil {
			log.Println("[DATABASE SEED] Berhasil membuat akun dummy: admin / 1234")
		} else {
			log.Println("[DATABASE SEED ERROR] Gagal membuat akun dummy:", err)
		}
	}

	// === SEED DATA SENSOR AWAL (Mencegah Warning Record Not Found di GORM) ===
	var sensorCount int64
	db.Model(&models.SensorLog{}).Count(&sensorCount)
	if sensorCount == 0 {
		defaultSensor := models.SensorLog{
			CahayaAtap:      80,
			DapurSuhu:       28.0,
			DapurKelembapan: 60.0,
			DapurFlame:      0,
			KamarSuhu:       27.0,
			KamarKelembapan: 55.0,
			TamuGerak:       false,
		}
		if err := db.Create(&defaultSensor).Error; err == nil {
			log.Println("[DATABASE SEED] Berhasil membuat log sensor awal default.")
		} else {
			log.Println("[DATABASE SEED ERROR] Gagal membuat log sensor awal:", err)
		}
	}
}
