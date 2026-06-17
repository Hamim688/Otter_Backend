package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// INI DIA BIANG KEROKNYA! Variabel global DB wajib huruf besar biar kedetect di file lain
var DB *gorm.DB

func ConnectDB() {
	// Ngerakit DSN dari file .env rahasia lu
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("[DATABASE] Gagal konek ke PostgreSQL: ", err)
	}

	DB = db
	fmt.Println("[DATABASE] Konek ke PostgreSQL Berhasil! 🐘")
}