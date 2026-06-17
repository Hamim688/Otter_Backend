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

}
