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
}