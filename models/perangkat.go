package models

type Perangkat struct {
	ID              uint `gorm:"primaryKey;default:1" json:"id"` // Selalu ID 1
	LampuKamar      bool `gorm:"default:false" json:"lampu_kamar"`
	LampuTamu       bool `gorm:"default:false" json:"lampu_tamu"`
	LampuKamarMandi bool `gorm:"default:false" json:"lampu_kamar_mandi"`
	LampuDapur      bool `gorm:"default:false" json:"lampu_dapur"`
	KipasKamar      bool `gorm:"default:false" json:"kipas_kamar"`
	KecepatanKipas  int  `gorm:"default:255" json:"kecepatan_kipas"`
	BuzzerAlrm      bool `gorm:"default:false" json:"buzzer_alrm"`
	LedMerahDapur   bool `gorm:"default:false" json:"led_merah_dapur"`
	KunciPintuRfid  bool `gorm:"default:true" json:"kunci_pintu_rfid"`
}