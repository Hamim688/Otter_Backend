package models

import (
	"log"

	"gorm.io/gorm"
)

// Seed mempopulasikan database dengan data awal jika tabel kosong
func Seed(db *gorm.DB) {
	// === 1. SEED DATA DUMMY AKUN ===
	var userCount int64
	db.Model(&User{}).Count(&userCount)
	if userCount == 0 {
		dummyUser := User{
			Username:    "admin",
			Password:    "admin123",
			DisplayName: "Mimah Dudim",
			Role:        "Administrator Rumah Pintar",
			AvatarURL:   "",
		}
		if err := db.Create(&dummyUser).Error; err == nil {
			log.Println("[DATABASE SEED] Berhasil membuat akun dummy: admin / admin123")
		} else {
			log.Println("[DATABASE SEED ERROR] Gagal membuat akun dummy:", err)
		}
	}

	// === 2. SEED DATA DUMMY RFID ===
	var rfidCount int64
	db.Model(&RfidCard{}).Count(&rfidCount)
	if rfidCount == 0 {
		dummyCards := []RfidCard{
			{UID: "A1B2C3D4", NamaPemilik: "Hamim", Status: "aktif"},
			{UID: "90E1F2A3", NamaPemilik: "Rafa Enrico", Status: "aktif"},
			{UID: "E2F3A4B5", NamaPemilik: "Tamu Asing", Status: "menunggu"},
		}
		for _, card := range dummyCards {
			if err := db.Create(&card).Error; err != nil {
				log.Printf("[DATABASE SEED ERROR] Gagal mempopulasi kartu %s: %v\n", card.UID, err)
			}
		}
		log.Println("[DATABASE SEED] Berhasil mempopulasikan kartu RFID dummy.")
	}
}
