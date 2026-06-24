package config

import (
	"fmt"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var MQTTClient mqtt.Client

func ConnectMQTT(messageHandler mqtt.MessageHandler) {
	opts := mqtt.NewClientOptions()
	
	// Tarik data broker dan kredensial dari .env
	opts.AddBroker(os.Getenv("MQTT_BROKER")) 
	opts.SetClientID(os.Getenv("MQTT_CLIENT_ID"))
	opts.SetUsername(os.Getenv("MQTT_USER"))
	opts.SetPassword(os.Getenv("MQTT_PASS"))
	
	opts.SetDefaultPublishHandler(messageHandler)

	opts.OnConnect = func(c mqtt.Client) {
		fmt.Println("[MQTT] Konek ke Broker Berhasil! 🚀")
		c.Subscribe("otter_smarthome/sensor", 0, nil)
		c.Subscribe("otter_smarthome/rfid_terdaftar/scan_terbaru", 0, nil)
		c.Subscribe("otter_smarthome/ai_alert", 0, nil)
		c.Subscribe("otter_smarthome/device_boot", 0, nil)
	}

	opts.OnConnectionLost = func(c mqtt.Client, err error) {
		fmt.Printf("[MQTT] Koneksi Putus: %v\n", err)
	}

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal("[MQTT] Gagal konek: ", token.Error())
	}
	MQTTClient = client
}