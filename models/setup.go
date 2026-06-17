package models

import (
	"log"
	"gorm.io/gorm"
)

// MigrateDB bakal dipanggil di main.go buat nyetak tabel
func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(
		&User{},
		&RfidCard{},
		&Notification{},
		&Perangkat{},
		&Otomatisasi{},
		&SensorLog{},
	)
	
	if err != nil {
		log.Fatal("[DATABASE] Gagal AutoMigrate tabel! Error: ", err)
	}
	
	log.Println("[DATABASE] Semua tabel berhasil di-migrate ke PostgreSQL! 🚀")

	// === SEED DATA DUMMY AKUN ===
	var count int64
	db.Model(&User{}).Count(&count)
	if count == 0 {
		dummyUser := User{
			Username:    "admin",
			Password:    "admin123",
			DisplayName: "Mimah Dudim",
			Role:        "Administrator Rumah Pintar",
			AvatarURL:   "",
		}
		if err := db.Create(&dummyUser).Error; err == nil {
			log.Println("[DATABASE SEED] Berhasil membuat akun dummy: admin / admin123")
		}
	}

	// === SEED DATA DUMMY RFID ===
	db.Model(&RfidCard{}).Count(&count)
	if count == 0 {
		dummyCards := []RfidCard{
			{UID: "A1B2C3D4", NamaPemilik: "Hamim", Status: "aktif"},
			{UID: "90E1F2A3", NamaPemilik: "Rafa Enrico", Status: "aktif"},
			{UID: "E2F3A4B5", NamaPemilik: "Tamu Asing", Status: "menunggu"},
		}
		for _, card := range dummyCards {
			db.Create(&card)
		}
		log.Println("[DATABASE SEED] Berhasil mempopulasikan kartu RFID dummy.")
	}
}