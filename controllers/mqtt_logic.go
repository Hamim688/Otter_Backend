package controllers

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// INI WAJIB HURUF M BESAR DI AWAL: MessagePubHandler
// Fungsi ini adalah otak tengah lu pas nerima data dari ESP32
var MessagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	payload := string(msg.Payload())

	// Print ke terminal biar lu tau kalau ESP32 lu ngirim data
	fmt.Printf("[MQTT MASUK] Topik: %s \n[PAYLOAD]: %s\n\n", topic, payload)

	// --- CONTOH LOGIKA ROUTING TOPIC ---
	if topic == "otter_smarthome/rfid_terdaftar/scan_terbaru" {
		fmt.Println(">> Menerima scan UID RFID dari ESP32, siap diolah!")
		// Nanti logika nyocokin ke Database PostgreSQL ditaruh di sini
	}

	if topic == "otter_smarthome/sensor" {
		// Nanti logika insert histori sensor atau cek kebakaran ditaruh di sini
	}
}