#include <WiFi.h>
#include <PubSubClient.h>
#include <ArduinoJson.h>

// ==========================================
// 1. KREDENSIAL JARINGAN
// ==========================================
const char* ssid        = "@wifi";
const char* password    = "876543210";
const char* mqtt_server = "10.83.175.86"; // IP Laptop Arch lu
const int   mqtt_port   = 1883;
const char* mqtt_user   = "esp32_kipas";  // Beda user/client ID biar nggak tabrakan
const char* mqtt_pass   = "otter_esp_123";

// ==========================================
// 2. MAPPING PIN (Khusus Kipas)
// ==========================================
#define PIN_KIPAS_KAMAR 27

// ==========================================
// 3. INISIALISASI KOMPONEN
// ==========================================
WiFiClient espClient;
PubSubClient client(espClient);

// ==========================================
// 4. FUNGSI SETUP JARINGAN (WIFI & MQTT)
// ==========================================
void setup_wifi() {
  delay(10);
  Serial.println();
  Serial.print("Connecting to ");
  Serial.println(ssid);

  WiFi.begin(ssid, password);

  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  Serial.println("");
  Serial.println("WiFi connected");
  Serial.println("IP address: ");
  Serial.println(WiFi.localIP());
}

void reconnect() {
  while (!client.connected()) {
    Serial.print("Attempting MQTT connection...");
    // Pastikan Client ID beda dari ESP utama
    if (client.connect("ESP32_Kipas_SmartHome", mqtt_user, mqtt_pass)) {
      Serial.println("connected");
      // Subscribe ke topik perangkat untuk mendengarkan perintah kipas
      client.subscribe("otter_smarthome/perangkat");
    } else {
      Serial.print("failed, rc=");
      Serial.print(client.state());
      Serial.println(" try again in 5 seconds");
      delay(5000);
    }
  }
}

// ==========================================
// 5. CALLBACK: NERIMA TITAH DARI GOLANG
// ==========================================
void callback(char* topic, byte* payload, unsigned int length) {
  String messageTemp;
  for (int i = 0; i < length; i++) {
    messageTemp += (char)payload[i];
  }
  Serial.print("Message arrived on topic: ");
  Serial.print(topic);
  Serial.print(". Message: ");
  Serial.println(messageTemp);

  if (String(topic) == "otter_smarthome/perangkat") {
    JsonDocument doc;
    DeserializationError error = deserializeJson(doc, messageTemp);
    if (error) {
      Serial.print(F("deserializeJson() failed: "));
      Serial.println(error.f_str());
      return;
    }

    // Eksekusi Kipas Kamar saja, abaikan perangkat lain
    bool kipasOn = doc["kipas_kamar"] ? true : false;
    int speedVal = doc["kecepatan_kipas"]; // Nilai 0 - 255

    if (kipasOn) {
      if (speedVal <= 0) speedVal = 255; // Fallback: paksa full speed jika nilai 0 atau negatif
      ledcWrite(PIN_KIPAS_KAMAR, speedVal);
      Serial.print("Kipas NYALA, Speed (PWM): ");
      Serial.println(speedVal);
      
      // JIKA MENGGUNAKAN RELAY BIASA (Aktifkan baris di bawah dan matikan ledcWrite jika pakai relay)
      // digitalWrite(PIN_KIPAS_KAMAR, HIGH);
    } else {
      ledcWrite(PIN_KIPAS_KAMAR, 0);
      Serial.println("Kipas MATI.");
      
      // JIKA MENGGUNAKAN RELAY BIASA
      // digitalWrite(PIN_KIPAS_KAMAR, LOW);
    }
  }
}

// ==========================================
// 6. SETUP: DEKLARASI PIN & INISIALISASI
// ==========================================
void setup() {
  Serial.begin(115200);

  // Inisialisasi PWM Kipas Kamar (Native LEDC ESP32 Core v3.x)
  ledcAttach(PIN_KIPAS_KAMAR, 5000, 8); // Frekuensi 5000Hz, Resolusi 8-bit (0-255)
  ledcWrite(PIN_KIPAS_KAMAR, 0); // Set awal kipas mati

  // Mulai Jaringan
  setup_wifi();
  client.setServer(mqtt_server, mqtt_port);
  client.setCallback(callback);
}

// ==========================================
// 7. LOOP: MURNI SUBSCRIBE (NON-BLOCKING)
// ==========================================
void loop() {
  if (!client.connected()) reconnect();
  client.loop();
}
