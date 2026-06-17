package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	// SESUAIKAN DENGAN NAMA MODULE LU
	"Backend/config"
	"Backend/controllers"
	"Backend/models"
	"Backend/routes"
	"Backend/seeder"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println(">>> STARTING OTER BACKEND (SECURE MODE) <<<")

	// 0. LOAD FILE .env PERTAMA KALI
	err := godotenv.Load()
	if err != nil {
		log.Println("[WARNING] File .env tidak ditemukan! Menggunakan environment bawaan OS.")
	} else {
		fmt.Println("[SYSTEM] File .env berhasil dimuat! 🔒")
	}

	// 1. Inisialisasi Database
	config.ConnectDB()
	models.MigrateDB(config.DB)
	seeder.Seed(config.DB)

	// 2. Inisialisasi MQTT Broker
	config.ConnectMQTT(controllers.MessagePubHandler)

	// 3. Inisialisasi Fiber API Web Server
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})

	routes.SetupRoutes(app)

	// Ambil port dari .env, kalau kosong pakai default 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	go func() {
		fmt.Printf("[FIBER] Server API jalan di http://localhost:%s 🌐\n", port)
		if err := app.Listen(":" + port); err != nil {
			log.Panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	fmt.Println("\n[SYSTEM] Mematikan Backend OTER...")
	config.MQTTClient.Disconnect(250)
	fmt.Println("[SYSTEM] Backend Mati dengan aman. Sampai jumpa!")
}